# Changelog

## [1.3.0](https://github.com/AksharP5/Patchline/compare/v1.2.0...v1.3.0) (2026-01-12)


### Features

* add cache discovery and list scaffolding ([e0a9247](https://github.com/AksharP5/Patchline/commit/e0a92473f77f1288e1195df588363c3187d6e919))
* add cache invalidation ([b0f9006](https://github.com/AksharP5/Patchline/commit/b0f9006b43e67627ce7e2e440105c91edd7913e5))
* add npm registry helpers ([5a33b33](https://github.com/AksharP5/Patchline/commit/5a33b33a2c31c57a9d782f3ed7a38e871ed1f4ac))
* add outdated command ([abe0384](https://github.com/AksharP5/Patchline/commit/abe03848e44fe0a913daf0f72c3cf9566b9fe363))
* add snapshot and rollback commands ([14bb04a](https://github.com/AksharP5/Patchline/commit/14bb04ae14ab5a40d4c72697470f4b9b750d635a))
* add snapshot store ([34c540e](https://github.com/AksharP5/Patchline/commit/34c540eef8385e8d7e54b803cd9660a6bf0a6ee2))
* add sync command ([9b98959](https://github.com/AksharP5/Patchline/commit/9b9895954758cd7e187d6c6e8eb4d6e626068045))
* add upgrade command ([15f7553](https://github.com/AksharP5/Patchline/commit/15f755365eb0bb8cc327a7b30a8b31537e0d7db8))
* add version selection helpers ([5e47e93](https://github.com/AksharP5/Patchline/commit/5e47e9302d5b0db411fbb4c84308967b2a8f4bcb))
* align local plugin discovery ([87b174f](https://github.com/AksharP5/Patchline/commit/87b174f520a53b9a25c13d0e0d50d1b57f5fb23d))
* align plugin dir and Bun cache discovery ([f542af5](https://github.com/AksharP5/Patchline/commit/f542af53aa21c0e7e3be9de18226b5b8d52d9200))
* detect node_modules cache ([f0a203c](https://github.com/AksharP5/Patchline/commit/f0a203c4f9513e867c6abbc35bb93cb1c518bcb6))
* implement upgrade command ([7e25280](https://github.com/AksharP5/Patchline/commit/7e252800b1fad9e3d818535011d11e4addeacf0c))
* prefer custom config for upgrades ([d979f6b](https://github.com/AksharP5/Patchline/commit/d979f6b3392aac7a473cd88e863dc20a4884de82))
* support .opencode.json ([eb098db](https://github.com/AksharP5/Patchline/commit/eb098db85ffce5f1873a484e95f7c93d4d6ba1ab))
* support custom config env and JSONC parsing ([af8b0c2](https://github.com/AksharP5/Patchline/commit/af8b0c29504b200d707a4808f070b6cb9dc0695c))
* support opencode env overrides ([a2c51f6](https://github.com/AksharP5/Patchline/commit/a2c51f62d8820fea0bf98761c26849165a1ee870))
* sync command cache invalidation ([d3ac7ce](https://github.com/AksharP5/Patchline/commit/d3ac7cec838e95de0235961d99b7534d74ea5b41))
* update opencode plugin specs ([b4d3165](https://github.com/AksharP5/Patchline/commit/b4d31658f18b75fd9c6061e61f6dba4884962cdf))


### Bug Fixes

* accept jsonc trailing commas ([447d006](https://github.com/AksharP5/Patchline/commit/447d0063c1ab749979a3d5f346718e9ddce0e629))
* add context to cache operations ([5da1eb5](https://github.com/AksharP5/Patchline/commit/5da1eb5dd79cefcc717b6252f5d73617a3d06a02))
* clarify sync refresh results ([4e815fe](https://github.com/AksharP5/Patchline/commit/4e815fe33cb1fc49516ee9775cd400246953ce4f))
* document snapshot paths ([fa3d7c2](https://github.com/AksharP5/Patchline/commit/fa3d7c29425d5d0f2b28ba20d6fb108becb4e6e4))
* fall back to installed base ([d5c8026](https://github.com/AksharP5/Patchline/commit/d5c80269c3b8b0304ffd71f8573e4bf21b1c171e))
* guard cache invalidation ([39b506d](https://github.com/AksharP5/Patchline/commit/39b506d0df3838ae2231c9fa43a9351297e44080))
* harden outdated registry checks ([36e11d6](https://github.com/AksharP5/Patchline/commit/36e11d69d0feeeb0209b05ca7a454add8b1ca302))
* keep legacy config paths ([37878d3](https://github.com/AksharP5/Patchline/commit/37878d33014d91157edf43f40c318433d800e784))
* release files ([566ee09](https://github.com/AksharP5/Patchline/commit/566ee09f7971cb7771021ca87a5253b3ba64a535))
* restrict cache paths ([2c7f61b](https://github.com/AksharP5/Patchline/commit/2c7f61b1cfc94ec2251ee7ed67badb0c293626b5))
* skip redundant upgrades ([2de558b](https://github.com/AksharP5/Patchline/commit/2de558b473cb3bd4f567323225af99ac8879461e))
* use documented plugin paths ([d96f0b2](https://github.com/AksharP5/Patchline/commit/d96f0b223d9f73194415729819da60b214390163))
* use vX.Y.Z release tags ([767bf0a](https://github.com/AksharP5/Patchline/commit/767bf0adf90039a61520470949f9c1de92b86e22))


### Documentation

* update README ([d9bf179](https://github.com/AksharP5/Patchline/commit/d9bf179ff1493d85a4b06f9956321b8e6f7e50f9))


### Chores

* add gitignore ([705033c](https://github.com/AksharP5/Patchline/commit/705033cc566401d7116fb218629a0c640bebe748))
* align module path and add base scaffolding ([ecc090a](https://github.com/AksharP5/Patchline/commit/ecc090a493b2e4a214c73ac65599dfd48df3f7a7))
* align module path and add base scaffolding ([b4d143d](https://github.com/AksharP5/Patchline/commit/b4d143d29ab033f4a73571c469f9bf7df39ee26b))
* **main:** release patchline 1.1.0 ([f9ae109](https://github.com/AksharP5/Patchline/commit/f9ae10957ff24335c1b0e971afa3e89b08ef4124))
* **main:** release patchline 1.1.0 ([bf4d707](https://github.com/AksharP5/Patchline/commit/bf4d7070d303716854789be72867faa3f32330a3))
* **main:** release patchline 1.2.0 ([14356d0](https://github.com/AksharP5/Patchline/commit/14356d07e875e1c875289a10db8cd4c394943fd8))
* **main:** release patchline 1.2.0 ([38d0a77](https://github.com/AksharP5/Patchline/commit/38d0a777e17cb93a2b91824a702fc8278d5de775))


### Refactors

* allow registry client injection ([c581453](https://github.com/AksharP5/Patchline/commit/c581453319a8d404aa534d9d6cb9644198734742))


### Tests

* add semver comparison coverage ([ae7aadb](https://github.com/AksharP5/Patchline/commit/ae7aadb75a612332353a45656a1a43a46bf82a77))
* add snapshot upgrade rollback flow ([03f1676](https://github.com/AksharP5/Patchline/commit/03f167692ab15eb23dc501ca96cfd109c1d3b19f))
* cover discovery and path candidates ([3c800c9](https://github.com/AksharP5/Patchline/commit/3c800c937d2ed7f040da880f15689e9f67a8a1a0))
* cover list/outdated logic ([72fb9e7](https://github.com/AksharP5/Patchline/commit/72fb9e7260bac1fe88b053087f7e4e5449941ef9))
* cover npm registry fetch ([0bbe0d9](https://github.com/AksharP5/Patchline/commit/0bbe0d90058a62fb852be7cbf00ec7e575eccf7a))
* harden registry/outdated cases ([268b553](https://github.com/AksharP5/Patchline/commit/268b553a9d189a3eb3b9470e08f97580076fd833))
* prefer opencode.json globally ([d8e4983](https://github.com/AksharP5/Patchline/commit/d8e4983d1e98e8acf465a8666c24e234a0827423))
* tighten flow assertions ([d5f1643](https://github.com/AksharP5/Patchline/commit/d5f16433013a5313fc1755514849318d0286a338))

## [1.2.0](https://github.com/AksharP5/Patchline/compare/patchline-v1.1.8...patchline-v1.2.0) (2026-01-12)


### Features

* add cache discovery and list scaffolding ([e0a9247](https://github.com/AksharP5/Patchline/commit/e0a92473f77f1288e1195df588363c3187d6e919))
* add cache invalidation ([b0f9006](https://github.com/AksharP5/Patchline/commit/b0f9006b43e67627ce7e2e440105c91edd7913e5))
* add npm registry helpers ([5a33b33](https://github.com/AksharP5/Patchline/commit/5a33b33a2c31c57a9d782f3ed7a38e871ed1f4ac))
* add outdated command ([abe0384](https://github.com/AksharP5/Patchline/commit/abe03848e44fe0a913daf0f72c3cf9566b9fe363))
* add snapshot and rollback commands ([14bb04a](https://github.com/AksharP5/Patchline/commit/14bb04ae14ab5a40d4c72697470f4b9b750d635a))
* add snapshot store ([34c540e](https://github.com/AksharP5/Patchline/commit/34c540eef8385e8d7e54b803cd9660a6bf0a6ee2))
* add sync command ([9b98959](https://github.com/AksharP5/Patchline/commit/9b9895954758cd7e187d6c6e8eb4d6e626068045))
* add upgrade command ([15f7553](https://github.com/AksharP5/Patchline/commit/15f755365eb0bb8cc327a7b30a8b31537e0d7db8))
* add version selection helpers ([5e47e93](https://github.com/AksharP5/Patchline/commit/5e47e9302d5b0db411fbb4c84308967b2a8f4bcb))
* align local plugin discovery ([87b174f](https://github.com/AksharP5/Patchline/commit/87b174f520a53b9a25c13d0e0d50d1b57f5fb23d))
* align plugin dir and Bun cache discovery ([f542af5](https://github.com/AksharP5/Patchline/commit/f542af53aa21c0e7e3be9de18226b5b8d52d9200))
* detect node_modules cache ([f0a203c](https://github.com/AksharP5/Patchline/commit/f0a203c4f9513e867c6abbc35bb93cb1c518bcb6))
* implement upgrade command ([7e25280](https://github.com/AksharP5/Patchline/commit/7e252800b1fad9e3d818535011d11e4addeacf0c))
* prefer custom config for upgrades ([d979f6b](https://github.com/AksharP5/Patchline/commit/d979f6b3392aac7a473cd88e863dc20a4884de82))
* support .opencode.json ([eb098db](https://github.com/AksharP5/Patchline/commit/eb098db85ffce5f1873a484e95f7c93d4d6ba1ab))
* support custom config env and JSONC parsing ([af8b0c2](https://github.com/AksharP5/Patchline/commit/af8b0c29504b200d707a4808f070b6cb9dc0695c))
* support opencode env overrides ([a2c51f6](https://github.com/AksharP5/Patchline/commit/a2c51f62d8820fea0bf98761c26849165a1ee870))
* sync command cache invalidation ([d3ac7ce](https://github.com/AksharP5/Patchline/commit/d3ac7cec838e95de0235961d99b7534d74ea5b41))
* update opencode plugin specs ([b4d3165](https://github.com/AksharP5/Patchline/commit/b4d31658f18b75fd9c6061e61f6dba4884962cdf))


### Bug Fixes

* accept jsonc trailing commas ([447d006](https://github.com/AksharP5/Patchline/commit/447d0063c1ab749979a3d5f346718e9ddce0e629))
* add context to cache operations ([5da1eb5](https://github.com/AksharP5/Patchline/commit/5da1eb5dd79cefcc717b6252f5d73617a3d06a02))
* clarify sync refresh results ([4e815fe](https://github.com/AksharP5/Patchline/commit/4e815fe33cb1fc49516ee9775cd400246953ce4f))
* document snapshot paths ([fa3d7c2](https://github.com/AksharP5/Patchline/commit/fa3d7c29425d5d0f2b28ba20d6fb108becb4e6e4))
* fall back to installed base ([d5c8026](https://github.com/AksharP5/Patchline/commit/d5c80269c3b8b0304ffd71f8573e4bf21b1c171e))
* guard cache invalidation ([39b506d](https://github.com/AksharP5/Patchline/commit/39b506d0df3838ae2231c9fa43a9351297e44080))
* harden outdated registry checks ([36e11d6](https://github.com/AksharP5/Patchline/commit/36e11d69d0feeeb0209b05ca7a454add8b1ca302))
* keep legacy config paths ([37878d3](https://github.com/AksharP5/Patchline/commit/37878d33014d91157edf43f40c318433d800e784))
* release files ([566ee09](https://github.com/AksharP5/Patchline/commit/566ee09f7971cb7771021ca87a5253b3ba64a535))
* restrict cache paths ([2c7f61b](https://github.com/AksharP5/Patchline/commit/2c7f61b1cfc94ec2251ee7ed67badb0c293626b5))
* skip redundant upgrades ([2de558b](https://github.com/AksharP5/Patchline/commit/2de558b473cb3bd4f567323225af99ac8879461e))
* use documented plugin paths ([d96f0b2](https://github.com/AksharP5/Patchline/commit/d96f0b223d9f73194415729819da60b214390163))


### Documentation

* update README ([d9bf179](https://github.com/AksharP5/Patchline/commit/d9bf179ff1493d85a4b06f9956321b8e6f7e50f9))


### Chores

* add gitignore ([705033c](https://github.com/AksharP5/Patchline/commit/705033cc566401d7116fb218629a0c640bebe748))
* align module path and add base scaffolding ([ecc090a](https://github.com/AksharP5/Patchline/commit/ecc090a493b2e4a214c73ac65599dfd48df3f7a7))
* align module path and add base scaffolding ([b4d143d](https://github.com/AksharP5/Patchline/commit/b4d143d29ab033f4a73571c469f9bf7df39ee26b))
* **main:** release patchline 1.1.0 ([f9ae109](https://github.com/AksharP5/Patchline/commit/f9ae10957ff24335c1b0e971afa3e89b08ef4124))
* **main:** release patchline 1.1.0 ([bf4d707](https://github.com/AksharP5/Patchline/commit/bf4d7070d303716854789be72867faa3f32330a3))


### Refactors

* allow registry client injection ([c581453](https://github.com/AksharP5/Patchline/commit/c581453319a8d404aa534d9d6cb9644198734742))


### Tests

* add semver comparison coverage ([ae7aadb](https://github.com/AksharP5/Patchline/commit/ae7aadb75a612332353a45656a1a43a46bf82a77))
* add snapshot upgrade rollback flow ([03f1676](https://github.com/AksharP5/Patchline/commit/03f167692ab15eb23dc501ca96cfd109c1d3b19f))
* cover discovery and path candidates ([3c800c9](https://github.com/AksharP5/Patchline/commit/3c800c937d2ed7f040da880f15689e9f67a8a1a0))
* cover list/outdated logic ([72fb9e7](https://github.com/AksharP5/Patchline/commit/72fb9e7260bac1fe88b053087f7e4e5449941ef9))
* cover npm registry fetch ([0bbe0d9](https://github.com/AksharP5/Patchline/commit/0bbe0d90058a62fb852be7cbf00ec7e575eccf7a))
* harden registry/outdated cases ([268b553](https://github.com/AksharP5/Patchline/commit/268b553a9d189a3eb3b9470e08f97580076fd833))
* prefer opencode.json globally ([d8e4983](https://github.com/AksharP5/Patchline/commit/d8e4983d1e98e8acf465a8666c24e234a0827423))
* tighten flow assertions ([d5f1643](https://github.com/AksharP5/Patchline/commit/d5f16433013a5313fc1755514849318d0286a338))

## [1.1.0](https://github.com/AksharP5/Patchline/compare/patchline-v1.0.0...patchline-v1.1.0) (2026-01-10)


### Features

* add cache discovery and list scaffolding ([e0a9247](https://github.com/AksharP5/Patchline/commit/e0a92473f77f1288e1195df588363c3187d6e919))
* add cache invalidation ([b0f9006](https://github.com/AksharP5/Patchline/commit/b0f9006b43e67627ce7e2e440105c91edd7913e5))
* add npm registry helpers ([5a33b33](https://github.com/AksharP5/Patchline/commit/5a33b33a2c31c57a9d782f3ed7a38e871ed1f4ac))
* add outdated command ([abe0384](https://github.com/AksharP5/Patchline/commit/abe03848e44fe0a913daf0f72c3cf9566b9fe363))
* add snapshot and rollback commands ([14bb04a](https://github.com/AksharP5/Patchline/commit/14bb04ae14ab5a40d4c72697470f4b9b750d635a))
* add snapshot store ([34c540e](https://github.com/AksharP5/Patchline/commit/34c540eef8385e8d7e54b803cd9660a6bf0a6ee2))
* add sync command ([9b98959](https://github.com/AksharP5/Patchline/commit/9b9895954758cd7e187d6c6e8eb4d6e626068045))
* add upgrade command ([15f7553](https://github.com/AksharP5/Patchline/commit/15f755365eb0bb8cc327a7b30a8b31537e0d7db8))
* add version selection helpers ([5e47e93](https://github.com/AksharP5/Patchline/commit/5e47e9302d5b0db411fbb4c84308967b2a8f4bcb))
* align local plugin discovery ([87b174f](https://github.com/AksharP5/Patchline/commit/87b174f520a53b9a25c13d0e0d50d1b57f5fb23d))
* align plugin dir and Bun cache discovery ([f542af5](https://github.com/AksharP5/Patchline/commit/f542af53aa21c0e7e3be9de18226b5b8d52d9200))
* detect node_modules cache ([f0a203c](https://github.com/AksharP5/Patchline/commit/f0a203c4f9513e867c6abbc35bb93cb1c518bcb6))
* implement upgrade command ([7e25280](https://github.com/AksharP5/Patchline/commit/7e252800b1fad9e3d818535011d11e4addeacf0c))
* prefer custom config for upgrades ([d979f6b](https://github.com/AksharP5/Patchline/commit/d979f6b3392aac7a473cd88e863dc20a4884de82))
* support .opencode.json ([eb098db](https://github.com/AksharP5/Patchline/commit/eb098db85ffce5f1873a484e95f7c93d4d6ba1ab))
* support custom config env and JSONC parsing ([af8b0c2](https://github.com/AksharP5/Patchline/commit/af8b0c29504b200d707a4808f070b6cb9dc0695c))
* support opencode env overrides ([a2c51f6](https://github.com/AksharP5/Patchline/commit/a2c51f62d8820fea0bf98761c26849165a1ee870))
* sync command cache invalidation ([d3ac7ce](https://github.com/AksharP5/Patchline/commit/d3ac7cec838e95de0235961d99b7534d74ea5b41))
* update opencode plugin specs ([b4d3165](https://github.com/AksharP5/Patchline/commit/b4d31658f18b75fd9c6061e61f6dba4884962cdf))


### Bug Fixes

* accept jsonc trailing commas ([447d006](https://github.com/AksharP5/Patchline/commit/447d0063c1ab749979a3d5f346718e9ddce0e629))
* add context to cache operations ([5da1eb5](https://github.com/AksharP5/Patchline/commit/5da1eb5dd79cefcc717b6252f5d73617a3d06a02))
* clarify sync refresh results ([4e815fe](https://github.com/AksharP5/Patchline/commit/4e815fe33cb1fc49516ee9775cd400246953ce4f))
* document snapshot paths ([fa3d7c2](https://github.com/AksharP5/Patchline/commit/fa3d7c29425d5d0f2b28ba20d6fb108becb4e6e4))
* fall back to installed base ([d5c8026](https://github.com/AksharP5/Patchline/commit/d5c80269c3b8b0304ffd71f8573e4bf21b1c171e))
* guard cache invalidation ([39b506d](https://github.com/AksharP5/Patchline/commit/39b506d0df3838ae2231c9fa43a9351297e44080))
* harden outdated registry checks ([36e11d6](https://github.com/AksharP5/Patchline/commit/36e11d69d0feeeb0209b05ca7a454add8b1ca302))
* keep legacy config paths ([37878d3](https://github.com/AksharP5/Patchline/commit/37878d33014d91157edf43f40c318433d800e784))
* restrict cache paths ([2c7f61b](https://github.com/AksharP5/Patchline/commit/2c7f61b1cfc94ec2251ee7ed67badb0c293626b5))
* skip redundant upgrades ([2de558b](https://github.com/AksharP5/Patchline/commit/2de558b473cb3bd4f567323225af99ac8879461e))
* use documented plugin paths ([d96f0b2](https://github.com/AksharP5/Patchline/commit/d96f0b223d9f73194415729819da60b214390163))


### Chores

* add gitignore ([705033c](https://github.com/AksharP5/Patchline/commit/705033cc566401d7116fb218629a0c640bebe748))
* align module path and add base scaffolding ([ecc090a](https://github.com/AksharP5/Patchline/commit/ecc090a493b2e4a214c73ac65599dfd48df3f7a7))
* align module path and add base scaffolding ([b4d143d](https://github.com/AksharP5/Patchline/commit/b4d143d29ab033f4a73571c469f9bf7df39ee26b))


### Refactors

* allow registry client injection ([c581453](https://github.com/AksharP5/Patchline/commit/c581453319a8d404aa534d9d6cb9644198734742))


### Tests

* add semver comparison coverage ([ae7aadb](https://github.com/AksharP5/Patchline/commit/ae7aadb75a612332353a45656a1a43a46bf82a77))
* add snapshot upgrade rollback flow ([03f1676](https://github.com/AksharP5/Patchline/commit/03f167692ab15eb23dc501ca96cfd109c1d3b19f))
* cover discovery and path candidates ([3c800c9](https://github.com/AksharP5/Patchline/commit/3c800c937d2ed7f040da880f15689e9f67a8a1a0))
* cover list/outdated logic ([72fb9e7](https://github.com/AksharP5/Patchline/commit/72fb9e7260bac1fe88b053087f7e4e5449941ef9))
* cover npm registry fetch ([0bbe0d9](https://github.com/AksharP5/Patchline/commit/0bbe0d90058a62fb852be7cbf00ec7e575eccf7a))
* harden registry/outdated cases ([268b553](https://github.com/AksharP5/Patchline/commit/268b553a9d189a3eb3b9470e08f97580076fd833))
* prefer opencode.json globally ([d8e4983](https://github.com/AksharP5/Patchline/commit/d8e4983d1e98e8acf465a8666c24e234a0827423))
* tighten flow assertions ([d5f1643](https://github.com/AksharP5/Patchline/commit/d5f16433013a5313fc1755514849318d0286a338))

## Changelog

All notable changes to this project will be documented in this file.
