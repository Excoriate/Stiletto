# Changelog

## [0.2.11](https://github.com/Excoriate/Stiletto/compare/v0.2.10...v0.2.11) (2023-05-22)


### Bug Fixes

* fix bug trying to fetch ecs task env vars from keys ([08e65ab](https://github.com/Excoriate/Stiletto/commit/08e65ab0e7aa6ce9fb062d7f0d147e855da170d8))

## [0.2.10](https://github.com/Excoriate/Stiletto/compare/v0.2.9...v0.2.10) (2023-05-22)


### Bug Fixes

* fix set-env functionality ([98a3842](https://github.com/Excoriate/Stiletto/commit/98a3842a8dfcc8873ca08dfc540ecc0f4a8926c3))

## [0.2.9](https://github.com/Excoriate/Stiletto/compare/v0.2.8...v0.2.9) (2023-05-22)


### Bug Fixes

* fix mechanism for obtaining slices and maps from cli options ([1b5a674](https://github.com/Excoriate/Stiletto/commit/1b5a6741833f9161943a5f4be26e44081ef6f961))

## [0.2.8](https://github.com/Excoriate/Stiletto/compare/v0.2.7...v0.2.8) (2023-05-22)


### Features

* add support for scanning env vars from prefixes, fix bug ([ad8c97a](https://github.com/Excoriate/Stiletto/commit/ad8c97a8c2edb3cc08a3367241471cc5dd64c43f))

## [0.2.7](https://github.com/Excoriate/Stiletto/compare/v0.2.6...v0.2.7) (2023-05-21)


### Bug Fixes

* add git repositories validations, add proper validations for terragrunt module target dir ([2cd6962](https://github.com/Excoriate/Stiletto/commit/2cd69625a859fd9e02470c6257846e2e787a8671))

## [0.2.6](https://github.com/Excoriate/Stiletto/compare/v0.2.5...v0.2.6) (2023-05-20)


### Features

* add support for env vars in ecs deploy command ([cd73e6c](https://github.com/Excoriate/Stiletto/commit/cd73e6c031c497d31656c3d3d5952cd606c6fd09))

## [0.2.5](https://github.com/Excoriate/Stiletto/compare/v0.2.4...v0.2.5) (2023-05-19)


### Bug Fixes

* resolve workdir in a smart manner, based on key dot or current dir ([6386662](https://github.com/Excoriate/Stiletto/commit/6386662bb59ee42e42857568d10b4f31abbb798b))

## [0.2.4](https://github.com/Excoriate/Stiletto/compare/v0.2.3...v0.2.4) (2023-05-19)


### Bug Fixes

* fix ecr command, fix .env wrong asked file when option isn't set ([a2aac80](https://github.com/Excoriate/Stiletto/commit/a2aac8066dad4144ee8a59dd01119ac34bd22649))

## [0.2.3](https://github.com/Excoriate/Stiletto/compare/v0.2.2...v0.2.3) (2023-05-18)


### Bug Fixes

* remove invalid condition when it's checking the pipeline mandatory settings ([1557f40](https://github.com/Excoriate/Stiletto/commit/1557f405b96e0c6f81b8c4b8bb9fd03f6fc83e32))

## [0.2.2](https://github.com/Excoriate/Stiletto/compare/v0.2.1...v0.2.2) (2023-05-17)


### Features

* add support for dotenv files, refactor CLI args from job init ([ad2a210](https://github.com/Excoriate/Stiletto/commit/ad2a21022cf912b554627af9334fa9f2f7fe2c09))


### Refactoring

* add command runner function ([506f38b](https://github.com/Excoriate/Stiletto/commit/506f38ba74310376ffa059bef41914e111bd563f))
* add missing tg commands on task ([d0cfa2e](https://github.com/Excoriate/Stiletto/commit/d0cfa2ea5f00caaab4d5116628284a55eeb7524e))
* Add shared validation function for validate viper key config ([3cd4953](https://github.com/Excoriate/Stiletto/commit/3cd4953622592f550c44ec7d52afb4855ba316ec))
* add tg command, working version ([01639e4](https://github.com/Excoriate/Stiletto/commit/01639e4cb3867ce3204acad9a6c600c834c24221))


### Other

* adjust release-please config ([016d8d9](https://github.com/Excoriate/Stiletto/commit/016d8d9f1262be40ea4f5da087e4e91522d81c6e))

## [0.2.1](https://github.com/Excoriate/Stiletto/compare/v0.2.0...v0.2.1) (2023-04-10)


### Features

* add first commit, basic functionality in place ([7764cf5](https://github.com/Excoriate/Stiletto/commit/7764cf5075263b169989e759e256503c913375f1))
* first commit ([340f071](https://github.com/Excoriate/Stiletto/commit/340f071aef0d669e491f405f07f60c3ccd10fc1b))
* first commit ([b2ab70c](https://github.com/Excoriate/Stiletto/commit/b2ab70c21a0aabda7f22df7c2da28a41b8820a40))
* first stable structure ([4e03e9a](https://github.com/Excoriate/Stiletto/commit/4e03e9ad2be6b4c9ca7f70274e33593673002497))


### Other

* add manifest version for release-please ([8e56013](https://github.com/Excoriate/Stiletto/commit/8e5601312744642c64420e16038b7184798fe371))
* add manifest version for release-please ([0fb32ac](https://github.com/Excoriate/Stiletto/commit/0fb32ac988b0a13e07837298d2ab0804474978d4))
* add publisher token for brew formula ([eee73db](https://github.com/Excoriate/Stiletto/commit/eee73dbb44df9aacee8659ea551415588f006b99))
