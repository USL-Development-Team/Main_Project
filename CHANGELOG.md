# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0](https://github.com/USL-Development-Team/Main_Project/compare/v1.0.1...v1.1.0) (2025-08-17)


### ✨ Features

* remove hardcoded placeholder domains and add configuration validation ([5689110](https://github.com/USL-Development-Team/Main_Project/commit/568911040804f5c49f6cbd70fcac159a02e456eb)), closes [#24](https://github.com/USL-Development-Team/Main_Project/issues/24)


### 🐛 Bug Fixes

* **oauth:** remove hardcoded placeholder domains in production config ([da79975](https://github.com/USL-Development-Team/Main_Project/commit/da799751da54220290a4c89774346a372b5c8caf)), closes [#24](https://github.com/USL-Development-Team/Main_Project/issues/24)


### ♻️ Code Refactoring

* improve code readability and organization ([0fb30c4](https://github.com/USL-Development-Team/Main_Project/commit/0fb30c44d055dbf8ee55b39e1194e6e652e4ea88))

## [1.0.1](https://github.com/USL-Development-Team/Main_Project/compare/v1.0.0...v1.0.1) (2025-08-17)


### 🐛 Bug Fixes

* replace hardcoded localhost oauth redirects with environment-aware urls ([431fb5d](https://github.com/USL-Development-Team/Main_Project/commit/431fb5dd16eeaddf44486787859a4305e36f24c4)), closes [#14](https://github.com/USL-Development-Team/Main_Project/issues/14)
* skip .env file loading in platform environments ([#12](https://github.com/USL-Development-Team/Main_Project/issues/12)) ([ae43255](https://github.com/USL-Development-Team/Main_Project/commit/ae43255f173d62327507b0b57f42b4c7842ff8bf)), closes [#11](https://github.com/USL-Development-Team/Main_Project/issues/11)
* use render deploy hook for deployments instead of missing api credentials ([07d3327](https://github.com/USL-Development-Team/Main_Project/commit/07d332769504a93895900da72efcd41f1915a2b2))


### ♻️ Code Refactoring

* complete oauth appbaseurl integration with centralized config ([f477bd3](https://github.com/USL-Development-Team/Main_Project/commit/f477bd3de1ba6170c3bac637e6d6421821004f25))
* improve code readability and maintainability ([c6a7509](https://github.com/USL-Development-Team/Main_Project/commit/c6a7509fddb1f77f585895ee4e08eaef20c61836))

## 1.0.0 (2025-08-17)


### ✨ Features

* add modern api patterns with pagination and filtering ([5438490](https://github.com/USL-Development-Team/Main_Project/commit/54384909ead9bc29fbcfab637953bef1d2846d83))
* add render blueprint for production deployment ([0ce2303](https://github.com/USL-Development-Team/Main_Project/commit/0ce2303f6d49f2ed3fc347b34f86dfbe0b4960fe))
* beautify codebase for improved readability and maintainability ([1105180](https://github.com/USL-Development-Team/Main_Project/commit/1105180215fb22a818346b7f46bcd66ecca0560d))
* clean up debug logs and obvious comments ([351659d](https://github.com/USL-Development-Team/Main_Project/commit/351659dd77c7be79f2b60cf7acb5012ceea08665))
* complete phase 2 template integration and fix linting issues ([f761863](https://github.com/USL-Development-Team/Main_Project/commit/f761863cd7a9ff75ad8120d3c74b1c2025923595))
* implement comprehensive validation system for usl tracker crud ([699b4e8](https://github.com/USL-Development-Team/Main_Project/commit/699b4e80b741d7e9c2c65561403605a7a3422c71))
* implement phase 4 client-side validation framework ([8ef4670](https://github.com/USL-Development-Team/Main_Project/commit/8ef4670ce6350505834d1f5c977c13bcc34fad42))
* initial go project setup with supabase integration ([865897b](https://github.com/USL-Development-Team/Main_Project/commit/865897b7d0a141747825b168d0edf40f8c56117a))
* initial production release with automation and templates ([#5](https://github.com/USL-Development-Team/Main_Project/issues/5)) ([d064c64](https://github.com/USL-Development-Team/Main_Project/commit/d064c64d3ae10daea97a248adcf5aaf96715ec77))
* integrate usl with trueskill calculation system ([4b316b7](https://github.com/USL-Development-Team/Main_Project/commit/4b316b776dceeee9d65484defefd22939829f155))
* remove obvious comments for cleaner code ([6d2eab6](https://github.com/USL-Development-Team/Main_Project/commit/6d2eab62a7ca3d9b08a710c11d3b2aef05460db8))
* start validation system refactoring session ([b0dbe59](https://github.com/USL-Development-Team/Main_Project/commit/b0dbe59fabe4bb65510d87cebce2dc75ca763cf6))


### 🐛 Bug Fixes

* remove problematic usl import script ([8661ae1](https://github.com/USL-Development-Team/Main_Project/commit/8661ae1d2bbdd60c2e03bb0dc9f66498fc9f6f0e))


### ♻️ Code Refactoring

* beautify code structure and focus test suite ([8cf94cc](https://github.com/USL-Development-Team/Main_Project/commit/8cf94cc44f9e0d08155f87e987b2d3d5ff86a729))
* beautify codebase with cleaner naming and structure ([4852bba](https://github.com/USL-Development-Team/Main_Project/commit/4852bbafd7e936f7b4eb70d7c4ac139c2a2a1b29))
* improve code readability and naming ([99742f8](https://github.com/USL-Development-Team/Main_Project/commit/99742f8839f42d29f7dd437bdedff23a49f855c8))
* improve code structure and performance ([9e82161](https://github.com/USL-Development-Team/Main_Project/commit/9e8216129b3467130e8ce24fd0e06f7b7ef80d70))
* remove obvious comments that restate code ([636a720](https://github.com/USL-Development-Team/Main_Project/commit/636a720e140ebb98da3de932a93f1e2fa6823128))
* remove redundant comments from go files ([07b58f4](https://github.com/USL-Development-Team/Main_Project/commit/07b58f4dd073d25ed406308f1929e0e1134313bd))
* rename environment files to follow proper pipeline naming ([cbbff7d](https://github.com/USL-Development-Team/Main_Project/commit/cbbff7dea4ce4d7d5a9fbfe0dc11fc4047004d76))
* **usl/handlers:** improve code quality and maintainability ([5ef3680](https://github.com/USL-Development-Team/Main_Project/commit/5ef36804089c2a364436acfc1ef014147a2b3b5f))


### 🧪 Tests

* add comprehensive validation safety and security tests ([97db658](https://github.com/USL-Development-Team/Main_Project/commit/97db658a88e581c40385f4576ce1fcbae7a6882d))
* add comprehensive validation tests and fix client-side validation ([c17ef87](https://github.com/USL-Development-Team/Main_Project/commit/c17ef87840d8eab03546fc4fd15c1aa7c6f572ce))
