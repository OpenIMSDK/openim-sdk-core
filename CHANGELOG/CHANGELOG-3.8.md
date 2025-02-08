## [v3.8.3](https://github.com/openimsdk/openim-sdk-core/releases/tag/v3.8.3) 	(2025-01-08)

### New Features
* feat: merge main bug fix to main. [#637](https://github.com/openimsdk/openim-sdk-core/pull/637)
* feat: improve merge in milestone and merged handle logic. [Created [#797](https://github.com/openimsdk/openim-sdk-core/pull/797)
* feat: add a function to quickly retrieve the context messages for a given message. [#828](https://github.com/openimsdk/openim-sdk-core/pull/828)

### Bug Fixes
* fix: create index failed when table name has `-`. [Created [#795](https://github.com/openimsdk/openim-sdk-core/pull/795)
* fix: change the table name and add escaping. [Created [#803](https://github.com/openimsdk/openim-sdk-core/pull/803)
* fix: change errs to custom errs avoid sdk panic. [Created [#802](https://github.com/openimsdk/openim-sdk-core/pull/802)
* fix: get reverse history message change. [#805](https://github.com/openimsdk/openim-sdk-core/pull/805)
* fix: add server isEnd determination criteria for message retrieval. [#812](https://github.com/openimsdk/openim-sdk-core/pull/812)
* fix: login user's info maybe empty when app reinstall. [#816](https://github.com/openimsdk/openim-sdk-core/pull/816)
* fix: search message do not filter voice message when keyword is empty. [#820](https://github.com/openimsdk/openim-sdk-core/pull/820)
* fix: quote message change to revoke message when app from background … [#826](https://github.com/openimsdk/openim-sdk-core/pull/826)

### Refactors
* refactor: add a parameter to locate messages and reverse pull messages to avoid UI data interference. [#833](https://github.com/openimsdk/openim-sdk-core/pull/833)
* refactor: remove fetch messages instead of search message clear cache. [#835](https://github.com/openimsdk/openim-sdk-core/pull/835)

### Builds
* build: update PR body. [Created [#798](https://github.com/openimsdk/openim-sdk-core/pull/798)

### Others
* bump: update go mod dependency version to latest. (#632) [#633](https://github.com/openimsdk/openim-sdk-core/pull/633)
* merge: update release-v3.8 with main changes [#735](https://github.com/openimsdk/openim-sdk-core/pull/735)

**Full Changelog**: [v3.8.2...v3.8.3](https://github.com/openimsdk/openim-sdk-core/compare/v3.8.2...v3.8.3)

