# Changelog

## [1.1.1](https://github.com/meysam81/oneoff/compare/v1.1.0...v1.1.1) (2025-11-30)


### Bug Fixes

* add cgo aware sqlite to migrations as well ([1a9db02](https://github.com/meysam81/oneoff/commit/1a9db024eaea015a155725bbef9683fb8dbd1981))
* **CI:** handle windows specific signals ([84d2687](https://github.com/meysam81/oneoff/commit/84d2687a946f420703249cffd1766d22f9443c19))

## [1.1.0](https://github.com/meysam81/oneoff/compare/v1.0.2...v1.1.0) (2025-11-30)


### Features

* **docs:** redesign landing page for higher conversions ([#36](https://github.com/meysam81/oneoff/issues/36)) ([8307b61](https://github.com/meysam81/oneoff/commit/8307b61eba1a2b48c72f858cd0e6177f38f204e6))
* **landing-page:** add comprehensive mobile responsiveness ([744b7e8](https://github.com/meysam81/oneoff/commit/744b7e8599d069d95110be4fa5e7a2895a10965c))
* plan first phase development roadmap ([#27](https://github.com/meysam81/oneoff/issues/27)) ([fd52dfc](https://github.com/meysam81/oneoff/commit/fd52dfcd72a8a0d9cf6414f8f3a223bd98ea87af))


### Bug Fixes

* add context to init and use cgo aware sqlite ([6065170](https://github.com/meysam81/oneoff/commit/6065170b5554b056b5ff264c1a74b99ed9a85c8a))
* add OS detection to landing page ([cc29ff2](https://github.com/meysam81/oneoff/commit/cc29ff2353c34e8df01c3dc5563e727fc33123e0))
* **CI:** sign the blobs ([b070e1a](https://github.com/meysam81/oneoff/commit/b070e1a89075870a809de95a0b6751d32f34d4f9))
* consolidate header and footer of landing page ([8102e46](https://github.com/meysam81/oneoff/commit/8102e46fc31064801c9a1ab722ed13abf23f3241))
* delegate more of OS tasks to ts file ([da49533](https://github.com/meysam81/oneoff/commit/da4953301448d1ae35e4477aa57482315812a5f0))
* detect OS and render landing page accordingly client-side ([e7be21e](https://github.com/meysam81/oneoff/commit/e7be21e44f6dbd6c7aadf93b1814957be44ddf24))
* Fetch GitHub version at build time for Hero ([#34](https://github.com/meysam81/oneoff/issues/34)) ([f559a08](https://github.com/meysam81/oneoff/commit/f559a08965e083aaaebbcdedafb3969f7cc2dafe))
* make features fully rectangular ([5da86b1](https://github.com/meysam81/oneoff/commit/5da86b1244be7ba9064c3b499ddfbb1d9c82f19c))
* run dev on all hosts ([21dad3f](https://github.com/meysam81/oneoff/commit/21dad3f4ef4ae9f7372d7d545e5a0a2f1571000e))
* use the header and footer from base layout in catalog ([99aea27](https://github.com/meysam81/oneoff/commit/99aea27e7a4fbe4ca0eabd433665863875f082e7))

## [1.0.2](https://github.com/meysam81/oneoff/compare/v1.0.1...v1.0.2) (2025-11-27)


### Bug Fixes

* **CI:** use official goreleaser action ([a9d3900](https://github.com/meysam81/oneoff/commit/a9d3900ddee78491c74682ad969d0754a125fee7))

## [1.0.1](https://github.com/meysam81/oneoff/compare/v1.0.0...v1.0.1) (2025-11-25)


### Bug Fixes

* add syft ([a1ba8be](https://github.com/meysam81/oneoff/commit/a1ba8be9a26047254efa4930109420a6587d8431))

## 1.0.0 (2025-11-25)


### Features

* add clone and immediate execution ([5e46487](https://github.com/meysam81/oneoff/commit/5e464871abdd13480f34521300217e1a26031ff4))
* add main entry point and fix build issues ([#4](https://github.com/meysam81/oneoff/issues/4)) ([8b866ef](https://github.com/meysam81/oneoff/commit/8b866ef40c7adc86170c1e78fceaadfa47da8d59))
* add SPA routing support with NotFound handler ([#5](https://github.com/meysam81/oneoff/issues/5)) ([fca2490](https://github.com/meysam81/oneoff/commit/fca2490d4aa92df242ba7c7db34e276e6fb537fc))
* build OneOff job scheduler application ([#1](https://github.com/meysam81/oneoff/issues/1)) ([5fd5a7f](https://github.com/meysam81/oneoff/commit/5fd5a7f89713ddf9af8c98747bef64bf34611678))
* change Docker job env vars input to .env format ([#19](https://github.com/meysam81/oneoff/issues/19)) ([ddc04d2](https://github.com/meysam81/oneoff/commit/ddc04d26f8852c0af3d084fec30b1bf03a33c70e))
* comprehensive frontend performance optimizations ([#14](https://github.com/meysam81/oneoff/issues/14)) ([51c96e2](https://github.com/meysam81/oneoff/commit/51c96e2b54f0133c1058a8a99609c6d35beb0dca))
* Set up CI pipeline with linting and goreleaser ([#6](https://github.com/meysam81/oneoff/issues/6)) ([0e367f2](https://github.com/meysam81/oneoff/commit/0e367f28bedbe4302334abd10568b49f1f38a1a9))


### Bug Fixes

* address Docker command button disappearing in job modal ([#21](https://github.com/meysam81/oneoff/issues/21)) ([981aadc](https://github.com/meysam81/oneoff/commit/981aadcda30f838f8049137d88cfb8657ea53d89))
* allow immediate execution as CTA ([045ae7f](https://github.com/meysam81/oneoff/commit/045ae7f372eadde94406a7825825ca56ae236db1))
* **CI:** disable cgo ([b603d54](https://github.com/meysam81/oneoff/commit/b603d549ef826d257ba2e362029dacecb41db210))
* **CI:** run release please independently ([db5e0c0](https://github.com/meysam81/oneoff/commit/db5e0c0843dad37273b31c54c054868b5e8af22b))
* **dev:** bring the main.go to the root ([9178b81](https://github.com/meysam81/oneoff/commit/9178b81bd1bcb91026970763e6b517e0c2b91776))
* do not track tmp ([8e286cf](https://github.com/meysam81/oneoff/commit/8e286cf04ad6b14cc83f569f2a3e8df9bc863b8d))
* parse datetime in go correctly ([6d65935](https://github.com/meysam81/oneoff/commit/6d659351c0384b5407d16820aa52b6ee0e901f8a))
* use vite compression2 ([b11ab85](https://github.com/meysam81/oneoff/commit/b11ab85df1cf7c64cb083031eb424bc3b68c9d92))
