# AUDIT — go-url-shortener

**Date :** 2026-05-28
**Périmètre :** commit `a86987d` (branche `main`)
**Stack :** Go 1.26 (Fiber + GORM + MySQL) · Quasar/Vue 3 + TypeScript
**Méthode :** revue manuelle exhaustive du dépôt (hors `node_modules`, `.git`, `dist`).

**Légende sévérité :**
- 🔴 **Critique** — exploitable directement, perte de données / compromission
- 🟠 **Élevé** — faille importante, contournable mais réaliste
- 🟡 **Moyen** — mauvaise pratique, dette technique notable
- 🔵 **Faible** — cosmétique / amélioration mineure

---

## 1. Architecture & organisation

### ✅ Points positifs
- Séparation classique **handler / repository / model** lisible (`server/handlers/`, `server/repositories/`, `server/models/`). Les repositories n'importent pas Fiber.
- Bootstrap commun centralisé : `cli/root.go:35` → `initConfigLoggerDatabase(initLogger, initDatabase bool)` permet aux sous-commandes Cobra de sauter le logger ou la DB.
- Cobra correctement utilisé pour `run` / `register` / `logs` (`cli/server.go`, `cli/user.go`, `cli/logs.go`).
- Routing Fiber idiomatique avec groupes (`server/routes.go:51`).

### ❌ Points négatifs
- 🟠 **404 inaccessible.** Le handler de 404 (`server/server.go:52-57`) est monté **après** `initJWT` (`server/server.go:47`). Toute URL inconnue renvoie un 401 JWT au lieu d'un 404. Le commentaire dans le code l'admet explicitement.
- 🟠 **Pas de RBAC.** `POST /api/v1/register` est dans `registerProtectedAPIRoutes` (`routes.go:64`), mais aucune notion de rôle ou d'admin n'existe : **tout utilisateur authentifié peut créer, lister, supprimer n'importe quel compte**. Privilege escalation triviale pour une application d'admin.
- 🟡 **Pas de struct de configuration.** Des `viper.GetString(...)` sont éparpillés partout (`server.go:79,213`, `cli/server.go:32`, `handlers/user.go:72,86`, `repositories/link.go:78`...). Difficile à tester / valider.
- 🟡 **Couplage handler → repository sans interface.** Aucun mock possible : tester un handler nécessite une vraie DB. Aucun test de handler n'existe.
- 🔵 Versions désynchronisées : `version = "1.3.0"` (`cli/root.go:14`), `CHANGELOG.md` (1.3.1), `client/package.json` (1.5.0).
- 🔵 Fichier mal nommé : `server/utils/valadator.go` (faute, devrait être `validator.go`).

---

## 2. Sécurité

### ✅ Points positifs
- Middlewares Fiber : `recover`, `requestid`, `limiter` configurable (`server/server.go:179-205`).
- BasicAuth devant `/doc/api-v1` et `/debug/pprof` (`routes.go:30`, `server.go:221`).
- Toutes les requêtes GORM utilisent des paramètres préparés → **pas de SQL injection**.
- Tokens reset password générés via `google/uuid` (v4 cryptographique, `handlers/user.go:250`).
- `Password` du modèle marqué `json:"-"` (`models/user.go:13`).

### ❌ Points négatifs

#### 🔴 Critiques
- **Hash mot de passe SHA-512 non salé.** `repositories/user.go:19,46,84,107` et `handlers/user.go:323` utilisent un SHA-512 nu. Cassable en GPU/rainbow tables en quelques secondes pour tout mdp < 12 caractères. La politique impose seulement `min=8` (`models/user.go:13`). **Tous les hashes en base doivent être considérés comme compromis.** `golang.org/x/crypto v0.52.0` (`go.mod:57`) est déjà tiré transitivement → migrer vers `bcrypt` (ou `argon2id`).
- **Secret JWT faible et versionné.** `JWT_SECRET=mySecretKeyForJWT` dans `.env.dist:34` *et* dans `.env:34` présent dans la copie de travail (avec `DB_PASSWORD=root`, `SERVER_BASICAUTH_PASSWORD=toto`). Vérifier `git log -- server/.env` pour confirmer qu'il n'est pas en historique.
- **Énumération d'utilisateurs.** `ForgottenPassword` (`handlers/user.go:240-245`) renvoie un 404 explicite si l'email n'existe pas. Toujours renvoyer 200/204.

#### 🟠 Élevés
- **Login vulnérable au timing & énumération.** `repositories/user.go:22` fait `Where(&User{Username: u, Password: hash})` — pas de comparaison constant-time, la latence diffère selon l'existence de l'utilisateur. `handlers/user.go:56-65` distingue `ErrRecordNotFound` (401) de l'autre erreur (500).
- **Bug SQL silencieux.** `repositories/user.go:140-152` : `SELECT u.password AS passwors` (faute). Le champ `Password` du struct anonyme reste vide, donc la vérification anti-réutilisation (`handlers/user.go:324`) ne marche jamais.
- **Token reset jamais invalidé en cas d'échec.** `UpdateUserPassword` n'appelle `DeletePasswordReset` que sur succès → un attaquant peut réessayer indéfiniment.
- **CORS dangereux par défaut.** `.env.dist:38` : `CORS_ALLOW_ORIGINS=*` + `CORS_ALLOW_CREDENTIALS=true` (`.env.dist:41`). De plus, la concaténation `strings.Join(..., ", ")` (`server.go:159`) ajoute des espaces qu'Fiber n'attend pas.
- **Open redirect.** `RedirectURL` (`handlers/link.go:179-191`) redirige sur `link.URL` sans valider le schéma. `models/link.go:21` n'impose que `validate:"required"`. Un compte compromis peut produire des phishing-links sur le domaine du raccourcisseur.
- **Content-Type spoofable.** `UploadLink` (`handlers/link.go:200`) lit `file.Header["Content-Type"][0]` — header client-side, **panic potentiel** si absent (`recover` sauve, mais DoS trivial), aucune vérification d'extension/contenu.
- **Aucun en-tête de sécurité HTTP.** Pas de CSP, ni `X-Frame-Options`, ni `X-Content-Type-Options`, ni HSTS, ni `Referrer-Policy`.
- **pprof exposé en production.** `SERVER_PPROF=true` par défaut (`.env.dist:26`), protégé par BasicAuth `toto:toto`. `pprof` expose heap & goroutines → fuite potentielle du secret JWT en mémoire.
- **Pas de révocation JWT.** Le "logout" client (`stores/user.ts:49-53`) efface juste sessionStorage. Le token reste valide côté serveur jusqu'à expiration (24h par défaut).
- **Aucune protection bruteforce login.** Le rate-limit est global, pas spécifique à `/login`.

#### 🟡 Moyens
- **PII dans le JWT.** Les claims contiennent `username, lastname, firstname, createdAt` (`handlers/user.go:77-80`) — réversibles côté client. Seul `id` suffirait. Absence de `iss`, `aud`, `sub`.
- **Asymétrie d'algo JWT.** Création figée sur `HS512` (`handlers/user.go:68`), vérification dépend de `JWT_ALGO` (`server.go:229`). Changer la conf casse silencieusement la signature.
- **Mots de passe logués en clair.** `server.go:108-118` log `zap.ByteString("body", c.Body())` sur les erreurs → fuite des payloads `/login`, `/register`, `/update-password`.
- **Secret JWT lu à chaque login** (`handlers/user.go:86`), sans contrôle de non-vacuité.
- **Validation faible des URLs.** `LinkForm` (`models/link.go:20-24`) n'a pas de `url` validator ni de `max=`.
- **Pas d'unique constraint sur `Link.URL`** → le compteur "trop de doublons" (`repositories/link.go:70-74`) tolère un état dégradé.

#### 🔵 Faibles
- Le `register` CLI affiche le mot de passe en clair dans la console (`cli/user.go:91`) — risque dans `history` / logs shell.
- `nbf` positionné au même instant que `iat` (`handlers/user.go:82-83`) — inutile.

---

## 3. Backend Go — qualité

### ✅ Points positifs
- Structure modulaire claire, packages cohérents.
- Logger `go.uber.org/zap` (performant, structuré).
- Dépendances à jour : `fiber/v2 v2.52.13`, `gorm v1.31.1`, `golang-jwt/jwt/v5 v5.3.1`.
- Tests via `stretchr/testify` (`db/database_test.go`, `logger_test.go`, `utils/url_shortener_test.go`).
- Bench sur `GenerateShortLink`.
- Shutdown propre sur SIGINT (`server.go:61-66`).

### ❌ Points négatifs
- 🟠 **Aucun `context.Context` propagé.** Tous les repositories prennent `*db.DB` direct sans `ctx` → impossibilité d'annuler une requête, pas de tracing distribué, pas de deadline. À combiner avec `db.WithContext(ctx)`.
- 🟠 **Erreurs ignorées :**
  - `db/database.go:89-92` : `db.DB()` → `err` jamais vérifié.
  - `db/database.go:99` : retour de `AutoMigrate` ignoré.
  - `server.go:71-74` : si `app.Listen` échoue, on continue.
  - `handlers/link.go:182-184` : `c.SendStatus(404)` sans `return`.
- 🟠 **Couverture de tests faible.** Aucun test sur Login, JWT, validators, hash, handlers, repositories, CLI. Pas de test d'intégration HTTP. `make test` (`Makefile:84`) n'utilise pas `-race` par défaut.
- 🟡 `GetAllLinks` (`repositories/link.go:22-52`) chaîne `q.Where(...)` en ignorant les valeurs de retour (GORM = comportement dépendant de la version).
- 🟡 `Count` ignore l'erreur retournée (`repositories/link.go:28`).
- 🟡 `GenerateShortLink` (`utils/url_shortener.go:33-45`) : tronque le hash SHA-256 à 64 bits (`big.Int.Uint64()`) → collisions possibles, aucun retry/check d'unicité à l'INSERT.
- 🟡 `Concurrency: 256 * 1024 * 1024` (`server.go:84`) : **256 millions de connexions concurrentes** — valeur absurde, probable copier-coller de `BodyLimit`.
- 🟡 `ReduceMemoryUsage: true` (`server.go:85`) et `UnescapePath: true` (`server.go:86`) — décisions non justifiées (compromis perf et risque de bypass).
- 🟡 `getGormLogOutput` (`db/database.go:126-143`) : fichier ouvert sans cleanup au shutdown.
- 🟡 Debug oublié : `fmt.Printf("ids=%v\n", links)` (`handlers/link.go:167`).
- 🟡 **Dockerfile cassé** : `RUN go get github.com/markbates/pkger/cmd/pkger` (déprécié depuis Go 1.16 et `embed`), références à `projects.json` et `config-docker.toml` introuvables (`Dockerfile:23-34`).
- 🔵 Goroutine de log sans backpressure (`logger.go:128-130`) — perte au shutdown.

---

## 4. Frontend Vue / Quasar

### ✅ Points positifs
- TypeScript activé, Pinia, Vue Router 4, Quasar 2.
- Séparation `api/ models/ stores/ pages/ components/ layouts/ services/ boot/`.
- Garde de route global (`router/index.ts:32-45`).
- Intercepteur Axios propre : injection `Authorization: Bearer`, redirect logout sur 401 (`boot/http.ts:93-129`).
- ESLint 9 + Biome (`package.json:12-14`).

### ❌ Points négatifs
- 🟠 **JWT stocké en `sessionStorage`** (`stores/user.ts:33`, `models/User.ts:38-43`). Accessible à tout JS du même origin → vulnérable XSS. Préférable : cookie `HttpOnly; Secure; SameSite=Strict` (refonte API requise) ou au minimum CSP stricte (absente).
- 🟠 **`User.toSession()` bugué.** `models/User.ts:37-39` : `JSON.stringify(User)` passe la **classe**, pas une instance. `fromSession()` cast `<User>JSON.parse(...)` → pas d'instance reconstituée. Heureusement jamais appelé directement.
- 🟠 **Aucun test frontend.** `"test": "echo \"No test specified\" && exit 0"` (`package.json:15`). Pas de Vitest, ni Cypress, ni Playwright.
- 🟠 **Domaine personnel codé en dur.** `quasar.config.mjs.dist:76` : `API_BASE_URL: 'https://apitic.fr/api/v1'`, idem `SORT_URL_BASE`. Devrait venir de variables d'env de build.
- 🟠 **URL factice.** `api/Link.ts:25,196` : `new URL('https://www.apitic.com/links')` utilisé pour extraire `pathname+search` — astuce trompeuse à supprimer.
- 🟠 **Instance Axios recréée par requête** (`boot/http.ts:89`) — surcoût et duplication d'intercepteurs.
- 🟡 `withCredentials = true` (`boot/http.ts:67`) couplé au CORS `*` est incohérent.
- 🟡 **Mise à jour user UI non implémentée** : `updateUser()` (`pages/users/List.vue:310-312`) retourne sans rien faire, alors que l'API existe.
- 🟡 `console.log(datetime)` oublié (`pages/users/List.vue:211`).
- 🟡 Beaucoup d'`any` (`boot/http.ts:14,17,18`, `pages/links/List.vue:462`) + `eslint-disable @typescript-eslint/no-explicit-any` → affaiblit la valeur de TypeScript.
- 🟡 Pas de gestion d'erreur globale (ErrorBoundary / handler).
- 🟡 `manage401Error` (`boot/http.ts:39-49`) traite à la fois `status === 401` *et* `data.code === 401` — code mort ou design dépassé.
- 🔵 `engines.node: "^26 || ^24 || ^22 || ^20 || ^18 || ^16"` (`package.json:42`) — Node 16 et 18 EOL.

---

## 5. Base de données / migrations

### ✅ Points positifs
- Modèles GORM clairs (tags `primaryKey`, `size`, `unique`, `index`).
- Soft-delete activé (`gorm.DeletedAt`).
- FK avec cascade : `PasswordReset.User` (`models/user.go:19`).
- Pool de connexions configurable (`db/database.go:90-92`).
- DSN avec `parseTime=True`.

### ❌ Points négatifs
- 🟠 **AutoMigrate uniquement** (`db/migration.go`). Pas de migrations versionnées (golang-migrate, goose, atlas). Aucun rollback, aucun audit du schéma, aucune suppression de colonnes obsolètes.
- 🟠 **Pas d'index unique sur `Link.URL`** alors que le code dépend de l'unicité (`repositories/link.go:70-74`). `gorm:"index;size:511"` est juste un index simple.
- 🟡 **Recherches non-sargables.** `url LIKE ? OR name LIKE ?` (`repositories/link.go:26,32,58`) → full scan à grand volume. Pas d'index FULLTEXT.
- 🟡 **MySQL 5.7 EOL** depuis octobre 2023 (`docker-compose.yml:19`).
- 🟡 `MYSQL_USER=root` + `MYSQL_RANDOM_ROOT_PASSWORD=yes` (`docker-compose.yml:26-28`) — conflit (`root` est créé par l'image).
- 🟡 `Password` indexé (`models/user.go:13`) — étrange et inutile.
- 🟡 `PasswordResets` : PK = `UserID` → un seul reset actif par user (acceptable, à expliciter).

---

## 6. DevOps / CI / observabilité

### ✅ Points positifs
- Dockerfile + docker-compose présents.
- Air pour live-reload (`.air.toml`).
- Logger structuré JSON via zap, configurable (stdout / file, niveaux).

### ❌ Points négatifs
- 🟠 **Aucune CI/CD.** Pas de `.github/workflows/`, pas de `.gitlab-ci.yml`. Aucun pipeline (lint, tests, govulncheck, trivy, dependabot).
- 🟠 **Dockerfile cassé** (cf. §3). Probablement plus utilisable en l'état.
- 🟠 **Pas de healthcheck.** `drill.yml` interroge `/health-check` qui **n'existe pas** dans `routes.go`. Indispensable pour Kubernetes / load balancers.
- 🟠 **Secrets versionnés** : `.env` présent dans le working tree (et probablement dans l'historique).
- 🟡 Pas de métriques Prometheus, pas de tracing OpenTelemetry.
- 🟡 Logs locaux uniquement (pas d'envoi syslog/ELK/Loki).
- 🟡 Port désynchronisé : `docker-compose.yml:8` mappe `9900:8888` mais `APP_PORT=3000` (`.env.dist:2`).
- 🟡 Pas de rotation des logs (`lumberjack` manquant).
- 🟡 Pas de `.dockerignore` → `.env` copié dans l'image.

---

## 7. Documentation

### ✅ Points positifs
- `CLAUDE.md` très complet (architecture, commandes, instructions).
- `ROUTES.md` et `ROUTES.http` (exemples curl/REST Client).
- `README.md` racine et `server/README.md`.
- `CHANGELOG.md` au format Keep-a-Changelog.

### ❌ Points négatifs
- 🟡 **Trois sources de version divergentes** : `cli/root.go` (1.3.0), `CHANGELOG.md` (1.3.1, 2024-03-20), `client/package.json` (1.5.0).
- 🟡 **`ROUTES.md` désynchronisé** : ne documente pas `forgotten-password`, `update-password`, `/links/upload`, `/links/export/csv`, `/links/selected`.
- 🟡 `TODO.md` réduit à 3 lignes, aucune roadmap.

---

## 8. Performance

### ✅ Points positifs
- Bench sur `GenerateShortLink`.
- Pool de connexions configurable.
- Export CSV en streaming via `SetBodyStreamWriter` (`handlers/link.go:293-324`).
- Cache 1h sur `/assets` (`routes.go:43`).

### ❌ Points négatifs
- 🟠 **Pas de cache du redirect.** `RedirectURL` (`handlers/link.go:179`) requête la DB à chaque hit. **Hotpath** d'un raccourcisseur — un cache Redis / in-memory (TTL court) est attendu.
- 🟠 **`Prefork: true` dangereux.** Fork du process → connection pool GORM non partagé, comportement imprévisible (`server.go:80`).
- 🟠 **N+1 latent.** `GetAllLinks` fait `Count` puis `Find` séparément (`repositories/link.go:24-50`), avec `LIKE %x%` → deux full scans.
- 🟡 `ExportCSVLinks` (`handlers/link.go:282`) sans limite de rows — risque mémoire/temps.
- 🟡 Pas de batch insert sur `UploadLink` (insertion ligne par ligne).
- 🟡 Pas de cache HTTP (ETag, `Cache-Control`) sur les routes API.

---

## 📊 Synthèse priorisée

### 🔴 À corriger immédiatement
1. **Migrer SHA-512 → bcrypt** (`repositories/user.go:19,46,84,107`, `handlers/user.go:323`).
2. **Régénérer `JWT_SECRET`** et purger `.env` du dépôt + vérifier `git log -- server/.env`.
3. **Énumération users** sur `ForgottenPassword` (`handlers/user.go:240`) → réponse identique.
4. **Activer un RBAC** sur `/register`, `/users/*`, `/links/*`.

### 🟠 À planifier rapidement
- Open redirect `/:id` (valider `http(s)://` + liste blanche optionnelle).
- Bug SQL `passwors` (`repositories/user.go:141`).
- 404 derrière JWT (réorganiser `server.go:47-57`).
- pprof + BasicAuth `toto:toto` ouverts en prod (`.env.dist:24-26`).
- En-têtes de sécurité (helmet, CSP, HSTS, X-Frame-Options).
- Redaction des bodies dans les logs (`server.go:114`).
- Migrations BD versionnées (golang-migrate / goose / atlas).
- Tests d'intégration + CI (GitHub Actions).
- Réparer le Dockerfile (supprimer pkger).
- JWT côté front : envisager cookie `HttpOnly`.

### 🟡 Améliorations recommandées
- Centraliser la config dans un struct typé + validation au démarrage.
- Propager `context.Context` jusqu'à GORM via `db.WithContext(ctx)`.
- Cache de redirection (Redis / `ristretto`).
- Aligner les versions (`cli/root.go`, `CHANGELOG.md`, `package.json`).
- Refondre la documentation `ROUTES.md` à partir d'OpenAPI.
- Tests Vitest + Playwright côté front.
- Observabilité : Prometheus + OpenTelemetry.
- Index FULLTEXT MySQL pour la recherche `LIKE %x%`.
- Migrer MySQL 5.7 → 8.x (ou MariaDB 11).

### 🔵 Petits chantiers
- Renommer `valadator.go` → `validator.go`.
- Supprimer `fmt.Printf("ids=...")` (`handlers/link.go:167`), `console.log(datetime)` (`pages/users/List.vue:211`), affichage mdp dans `register` (`cli/user.go:91`).
- Corriger `Concurrency: 256 * 1024 * 1024` (`server.go:84`).
- Aligner port docker-compose et `APP_PORT`.
- Retirer Node 16/18 de `engines` (`package.json:42`).
- Supprimer le domaine `apitic.fr` codé en dur (`quasar.config.mjs.dist`).

---

## 🛠️ Pistes d'amélioration (roadmap proposée)

### Phase 1 — Sécurisation (1-2 semaines)
- [ ] Migration `bcrypt` (cost ≥ 12) + script de re-hash transparent au prochain login (champ `password_algo` sur User).
- [ ] Sortir `.env` du dépôt, régénérer **tous** les secrets, ajouter `git-secrets` ou `gitleaks` en pre-commit.
- [ ] Implémenter un middleware RBAC (`role` sur User : `admin`, `user`) et appliquer à tous les endpoints sensibles.
- [ ] Helmet-like middleware Fiber : CSP, HSTS, X-Frame, Referrer-Policy.
- [ ] Rate-limit dédié `/login` + délai constant côté serveur (sleep jitter).
- [ ] Valider le schéma des URLs raccourcies (regex `^https?://` + liste noire optionnelle).
- [ ] Redaction des logs (zap field stripping pour `password`, `token`).

### Phase 2 — Robustesse (2-3 semaines)
- [ ] Migrations versionnées : choisir entre `golang-migrate`, `goose`, `atlas`. Geler `AutoMigrate` aux environnements dev.
- [ ] Struct `Config` typée + validation au démarrage (échec rapide si secret vide ou DSN absent).
- [ ] Propagation `context.Context` dans tous les repositories + `db.WithContext(c.UserContext())` dans les handlers.
- [ ] Interfaces sur les repositories pour mocker dans les tests handlers.
- [ ] Tests d'intégration (Fiber `app.Test()`) ciblant Login, CRUD User, CRUD Link, redirect, upload CSV.
- [ ] CI GitHub Actions : `go test -race -cover`, `golangci-lint`, `govulncheck`, `trivy fs`, build Docker.
- [ ] Healthcheck `/health` (liveness) et `/ready` (readiness avec ping DB).
- [ ] Réparer le Dockerfile (passer à `embed` natif), ajouter `.dockerignore`.

### Phase 3 — Frontend (1-2 semaines)
- [ ] Vitest + Testing Library sur les stores, services API et composants critiques.
- [ ] Playwright pour le golden path (login → créer link → redirect).
- [ ] Externaliser `API_BASE_URL` et `SORT_URL_BASE` en variables d'env de build (`import.meta.env`).
- [ ] Refactor `User` model : convertir en simple type/interface + helpers purs (pas de classe stockée en sessionStorage).
- [ ] Réfléchir au passage du JWT en cookie `HttpOnly` (impact CSRF → ajouter token CSRF double-submit).
- [ ] Supprimer les `any` restants, activer `strict: true` complet.

### Phase 4 — Performance & observabilité (1-2 semaines)
- [ ] Cache des redirects en mémoire (`ristretto`) ou Redis (TTL 5-15 min, invalidation à l'UPDATE/DELETE du link).
- [ ] Index FULLTEXT sur `Link.url` et `Link.name` + bascule des recherches `LIKE %x%` vers `MATCH ... AGAINST`.
- [ ] Endpoint `/metrics` Prometheus (`fiberprometheus`).
- [ ] Tracing OpenTelemetry (otelfiber + otelgorm).
- [ ] Rotation des logs avec `lumberjack`.

### Phase 5 — Polissage (1 semaine)
- [ ] Aligner les numéros de version (`cli/root.go`, `CHANGELOG.md`, `package.json`).
- [ ] Régénérer la doc API depuis une spec OpenAPI canonique.
- [ ] Désactiver `Prefork` ou le documenter explicitement comme incompatible.
- [ ] Tuning `Concurrency` à une valeur réaliste (`256 * 1024`).
- [ ] Migration MySQL 5.7 → MySQL 8 ou MariaDB 11.

---

## Verdict global

Le projet présente une **architecture saine et lisible** pour un raccourcisseur d'URL, avec une stack moderne et bien choisie. La séparation des couches est respectée, la documentation `CLAUDE.md` est exemplaire, et le code reste accessible à quelqu'un qui débarque.

**Mais l'audit révèle des failles de sécurité graves qui empêchent toute mise en production en l'état** : le hash SHA-512 des mots de passe et le secret JWT versionné sont des "showstoppers" absolus. L'absence totale de tests d'intégration, de CI/CD et de migrations versionnées trahit un projet encore en phase de prototype malgré une apparence de maturité.

Les bonnes nouvelles : tous les points critiques sont **corrigeables sans refonte architecturale**. Avec 4 à 6 semaines de travail focalisé sur la roadmap ci-dessus, le projet peut atteindre un niveau de qualité production-ready.
