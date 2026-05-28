# go-url-shortener

Service de raccourcissement d'URL avec une interface d'administration. Le backend Go génère et redirige des liens courts, une SPA Quasar/Vue gère les utilisateurs et les liens via une API JSON authentifiée.

## Fonctionnalités

- Génération de codes courts de 8 caractères (SHA-256 + base58)
- Redirection publique `GET /:id` vers l'URL d'origine
- Liens avec nom, date d'expiration et soft-delete (GORM)
- Authentification JWT, gestion des utilisateurs et mot de passe oublié
- Import CSV en masse, export CSV streamé de tous les liens filtrés
- Documentation API protégée par Basic Auth (`/doc/api-v1`)
- QR codes générés côté SPA pour partager un lien
- Logs structurés (zap) avec lecteurs CLI colorisés (`logs --server` / `--database`)

## Stack technique

### Backend (`server/`)

- **Langage** : Go 1.26
- **Framework HTTP** : Fiber v2
- **ORM** : GORM (driver MySQL)
- **CLI** : Cobra
- **Configuration** : Viper (fichier `.env`)
- **Logging** : zap (application), logger GORM personnalisé (SQL)
- **Auth** : `gofiber/contrib/jwt` + `golang-jwt/jwt/v5`
- **Validation** : `go-playground/validator/v10`
- **Templates** : `gofiber/template/html` (page de doc API)

### Frontend (`client/`)

- **Framework** : Vue 3 + Quasar 2 (Vite)
- **Langage** : TypeScript
- **State** : Pinia
- **Routing** : vue-router
- **HTTP** : axios
- **i18n** : vue-i18n
- **QR codes** : qrcode.vue
- **Lint** : ESLint 9 (flat config) + Biome 2
- **Format** : Biome

### Infrastructure

- MySQL 5.7 (via `docker-compose` ou instance externe)
- Dockerfile multi-stage pour le binaire serveur

## Architecture

### Vue d'ensemble

```
+---------------------+         +-------------------------+
|                     |         |                         |
|   Navigateur        |         |  Admin SPA (Quasar)     |
|  (lien court)       |         |  client/  port 9000     |
|                     |         |                         |
+----------+----------+         +------------+------------+
           |                                 |
           |  GET /:id                       |  /api/v1/*  (JWT)
           |                                 |
           v                                 v
+--------------------------------------------------------+
|                                                        |
|          Serveur Fiber (server/  port APP_PORT)        |
|                                                        |
|  +------------------+    +-------------------------+   |
|  |  Web routes      |    |  API routes  /api/v1    |   |
|  |  GET /:id        |    |  /login  /forgotten...  |   |
|  |  /doc/api-v1     |    |  -- JWT gate --         |   |
|  |  /assets         |    |  /users  /links         |   |
|  +--------+---------+    +-----------+-------------+   |
|           |                          |                 |
|           +-----------+--------------+                 |
|                       |                                |
|                       v                                |
|              +-----------------+                       |
|              |   handlers/     |  validation, JSON     |
|              +--------+--------+                       |
|                       |                                |
|                       v                                |
|              +-----------------+                       |
|              |  repositories/  |  GORM queries         |
|              +--------+--------+                       |
|                       |                                |
|                       v                                |
|              +-----------------+                       |
|              |   models/       |  AutoMigrate          |
|              +--------+--------+                       |
+-----------------------|--------------------------------+
                        |
                        v
                 +-------------+
                 |   MySQL     |
                 +-------------+
```

### Couches du serveur

Découpage standard handler → repository → model :

```
cmd/main.go
   |
   v
cli/  (Cobra: run | register | logs)
   |
   v
server.go  -- Fiber app, middlewares, JWT
   |
   v
routes.go  -- enregistre web + API publiques / protegees
   |
   +---> handlers/   (HTTP, validation via utils.ValidateStruct)
   |        |
   |        v
   +---> repositories/  (acces DB, retourne des models)
   |        |
   |        v
   +---> models/   (GORM structs, listes dans db/migration.go)
   |
   +---> db/    (connexion, logger GORM, AutoMigrate)
   +---> utils/ (generation code court base58, erreurs HTTP, validator)
```

Règles importantes :

- `repositories/` ne connaît pas Fiber.
- `handlers/` ne fait pas de SQL direct.
- Toute config est lue via `viper.Get*` (pas de struct centrale).
- Nouveaux modèles → ajouter dans `db/migration.go::modelsList` pour qu'AutoMigrate les prenne.

### Flux d'une requête type

```
Client SPA                  Fiber                 Repository           MySQL
   |                          |                       |                  |
   |--- POST /api/v1/links -->|                       |                  |
   |    + Bearer JWT          |                       |                  |
   |                          |-- valider JWT         |                  |
   |                          |-- ValidateStruct      |                  |
   |                          |-- GenerateShortLink   |                  |
   |                          |     (sha256+base58)   |                  |
   |                          |---- CreateLink ------>|                  |
   |                          |                       |---- INSERT ----->|
   |                          |                       |<---- OK ---------|
   |                          |<------ Link ----------|                  |
   |<------ 201 JSON ---------|                       |                  |
```

```
Visiteur                    Fiber                 Repository           MySQL
   |                          |                       |                  |
   |---- GET /:id ----------->|                       |                  |
   |                          |--- GetLinkByID ------>|                  |
   |                          |                       |---- SELECT ----->|
   |                          |                       |<---- row --------|
   |                          |<------ URL -----------|                  |
   |<--- 302 Location: URL ---|                       |                  |
```

### Génération du code court

`utils.GenerateShortLink` (server/utils/url_shortener.go) :

```
url + clé secrète + UnixNano  --SHA-256-->  hash 32 octets
                                    |
                                    v
                       BigInt(hash) en uint64
                                    |
                                    v
                          base58 (alphabet Bitcoin)
                                    |
                                    v
                       8 premiers caractères = ID
```

L'ajout d'`UnixNano` garantit l'unicité même pour la même URL ; en cas de collision sur 8 caractères, la contrainte de clé primaire fait échouer l'insert et la création peut être retentée.

## Démarrage rapide

### Prérequis

- Go 1.26
- Node.js 20+ et npm
- MySQL 5.7+ (local ou via Docker)

### Backend

```bash
cd server
cp .env.example .env       # à créer si absent — voir clés ci-dessous
make serve                 # ou: go run cmd/main.go run
```

Créer un utilisateur admin avant de pouvoir se connecter au front :

```bash
go run cmd/main.go register \
  -l Doe -f John -e john@example.com -p 'changeme!'
```

### Frontend

```bash
cd client
cp quasar.config.mjs.dist quasar.config.mjs   # ajuster API_BASE_URL si besoin
npm install
npm run dev                                   # http://localhost:9000
```

### Avec Docker

```bash
cd server
docker-compose up           # API exposée sur localhost:9900, MySQL en interne
```

## Variables d'environnement (`server/.env`)

| Clé | Description |
| --- | --- |
| `APP_NAME`, `APP_ENV`, `APP_ADDR`, `APP_PORT` | Identité et binding du serveur |
| `SERVER_PREFORK` | Active le mode prefork de Fiber |
| `SERVER_BASICAUTH_USERNAME` / `_PASSWORD` | Protège `/doc/api-v1` |
| `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`, `DB_CHARSET`, `DB_COLLATION`, `DB_LOCATION` | Connexion MySQL |
| `DB_MAX_IDLE_CONNS`, `DB_MAX_OPEN_CONNS`, `DB_CONN_MAX_LIFETIME` | Pool de connexions (durée en heures) |
| `DB_USE_AUTOMIGRATIONS` | Si `true`, lance GORM `AutoMigrate` au démarrage |
| `GORM_LOG_LEVEL`, `GORM_LOG_OUTPUT`, `GORM_LOG_FILE_PATH` | Logger SQL |

D'autres clés liées au JWT et aux middlewares sont lues directement dans `server.go` et `routes.go`.

## Commandes utiles

### Serveur (`cd server`)

| Commande | Effet |
| --- | --- |
| `make serve` | Démarre le serveur |
| `make serve-race` | Démarre avec `-race` |
| `make watch` | Live-reload via Air |
| `make build` | Compile `./go-url-shortener` |
| `make test` | `go test -cover ./...` |
| `make bench` | Benchmarks |
| `make logs` | Tail des logs serveur colorisés |
| `make view-cover-count` | Couverture HTML |
| `go-url-shortener register -l … -f … -e … -p …` | Créer un utilisateur admin |

### Client (`cd client`)

| Commande | Effet |
| --- | --- |
| `npm run dev` | Quasar dev server (HMR) |
| `npm run build` | Build production SPA dans `dist/` |
| `npm run lint` | ESLint + Biome lint |
| `npm run format` | Biome format `--write` |
| `npm run check` | Biome check + autofix (lint + format) |

## Documentation API

- `server/ROUTES.md` — liste exhaustive des endpoints avec exemples
- `server/ROUTES.http` — requêtes prêtes à l'emploi pour VS Code REST Client / JetBrains HTTP Client
- `GET /doc/api-v1` — doc en ligne (Basic Auth)

## Benchmark

```bash
cd server
drill --benchmark drill.yml --stats --quiet
```

## Structure du dépôt

```
.
├── CLAUDE.md            # Guide pour Claude Code
├── README.md
├── CHANGELOG.md
├── TODO.md
├── client/              # SPA Quasar/Vue 3
│   ├── src/
│   │   ├── api/         # Wrappers HTTP (LinkAPI, UserAPI)
│   │   ├── boot/        # http.ts (axios + JWT)
│   │   ├── components/
│   │   ├── i18n/
│   │   ├── layouts/
│   │   ├── models/
│   │   ├── pages/       # links, users, user (login, password)
│   │   ├── router/
│   │   ├── services/
│   │   └── stores/      # Pinia
│   ├── biome.json
│   ├── eslint.config.mjs
│   └── quasar.config.mjs.dist
└── server/              # API Go/Fiber
    ├── cmd/main.go      # entrée Cobra
    ├── cli/             # commandes run, register, logs
    ├── db/              # connexion, AutoMigrate
    ├── handlers/        # link.go, user.go, doc.go
    ├── repositories/    # link.go, user.go
    ├── models/          # link.go, user.go
    ├── utils/           # generation code court, validator, http_error
    ├── templates/       # pages HTML (doc API)
    ├── assets/          # statique servi sous /assets
    ├── routes.go
    ├── server.go
    ├── Makefile
    ├── Dockerfile
    └── docker-compose.yml
```

## Convention de commit

Les messages de commit suivent une structure simple : une ligne de résumé, puis une liste de puces taguées par emoji selon la nature du changement.

```
<résumé en une ligne, < 70 caractères, sans point final>

🔥 Feature: <ce qui a été ajouté>
♻️ Refactor: <ce qui a été restructuré>
🩹 Fix: <ce qui a été corrigé>
🚨 Test: <ce qui change côté tests>
📚 Doc: <ce qui change côté documentation>
🎨 Style: <formatage / style UI>
```

Règles :

- La première ligne est un résumé clair, sans emoji ni préfixe.
- Ne garder que les puces qui s'appliquent réellement au commit (pas de catégorie vide).
- Rester simple et lisible, y compris pour des non-natifs : phrases courtes, pas de jargon.
- Regrouper les fichiers liés sous une même puce plutôt qu'une puce par fichier.
- Ne pas inventer de changement absent du diff.

Exemple :

```
Switch client tooling to Biome and bump deps

🔥 Feature: add CLAUDE.md and README.md
♻️ Refactor: replace legacy ESLint config with flat config and Biome
🎨 Style: reformat client files with Biome (4-space indent, trailing commas)
```

## Licence

Voir `server/LICENSE`.
