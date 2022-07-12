### Changelog

All notable changes to this project will be documented in this file. Dates are displayed in UTC.

#### [v1.3.0](https://bitbucket.org/decimalteam/go-node/compare/v1.3.0..v1.2.15) - 2022-06-29

-  [`#84`](https://bitbucket.org/decimalteam/go-node/pull-requests/84) feature/32285-burn-coin
-  [`f6ee961`](https://bitbucket.org/decimalteam/go-node/commits/f6ee9619009060acf0ccfe7a9c49d8fb822e4a3e) feature/32285-burn-coin feat: added burning coins feature. Increased min supply for coins from 1 pip to 1 bip (10^18 pip).
-  [`446c9e7`](https://bitbucket.org/decimalteam/go-node/commits/446c9e7cf646c9f90dc60396c7b88dfc8770b133) up version to 1.3.0
-  [`15b5a17`](https://bitbucket.org/decimalteam/go-node/commits/15b5a17a499d8c67210496e1778c23bf78362ad2) feature/32285-burn-coin feat: added params to newly added error.

#### [v1.2.15](https://bitbucket.org/decimalteam/go-node/compare/v1.2.15..v1.2.14) - 2022-06-09

-  [`#83`](https://bitbucket.org/decimalteam/go-node/pull-requests/83) Jailing absent validators without slashing
-  [`961f43e`](https://bitbucket.org/decimalteam/go-node/commits/961f43ede874d69341fc5f117dbf54731db09a92) hotfix/31324 fix: test for new slash logic
-  [`6270d98`](https://bitbucket.org/decimalteam/go-node/commits/6270d9895af0965e67fe856244807414ba852282) hotfix/31324 fix: changed handling missed votes to jail absent validators in grace period too.
-  [`35bf09a`](https://bitbucket.org/decimalteam/go-node/commits/35bf09a5db4f09b8891c7cd2aaa311dce85290fe) up version to 1.2.15

#### [v1.2.14](https://bitbucket.org/decimalteam/go-node/compare/v1.2.14..v1.2.13) - 2022-06-02

-  [`#82`](https://bitbucket.org/decimalteam/go-node/pull-requests/82) Hotfix/28883 - Compensations
-  [`7a0ac35`](https://bitbucket.org/decimalteam/go-node/commits/7a0ac35f5cd6d12fd3ff2688b8b3864ff3f8457f) hotfix/28883 fix: supported compensation mechanism and specific function compensating wrong slashes happened at 9288729 block.
-  [`1cd5b38`](https://bitbucket.org/decimalteam/go-node/commits/1cd5b38ab5e78b2af5e61f1447ea03c573cd536a) hotfix/28883 fix: removed commas from amounts in NFT compensations.
-  [`03d3eb6`](https://bitbucket.org/decimalteam/go-node/commits/03d3eb6a1f84ff2331854fc34b8cea42db7a4ae5) hotfix/28883 fix: removed panics from compensation code to make it safer.

#### [v1.2.13](https://bitbucket.org/decimalteam/go-node/compare/v1.2.13..v1.2.12) - 2022-04-29

-  [`#81`](https://bitbucket.org/decimalteam/go-node/pull-requests/81) Preprod
-  [`#80`](https://bitbucket.org/decimalteam/go-node/pull-requests/80) Hotfix/27660
-  [`#78`](https://bitbucket.org/decimalteam/go-node/pull-requests/78) hotfix/27660 fix: rapid fix issue for case when all slots of validator are used.
-  [`70081fe`](https://bitbucket.org/decimalteam/go-node/commits/70081fe5782fc823224e65a30f334b4c9ed5b4fc) hotfix/27660 fix: add test
-  [`28b1bbf`](https://bitbucket.org/decimalteam/go-node/commits/28b1bbfbff8559a947e51b9b476f81772c620d4e) hotfix/27660 fix: add test
-  [`09e46e5`](https://bitbucket.org/decimalteam/go-node/commits/09e46e56c4d04bca1869c6e1960f2842e1bf2aa1) hotfix/27660 test: test with coins

#### [v1.2.12](https://bitbucket.org/decimalteam/go-node/compare/v1.2.12..v1.2.10) - 2022-03-30

-  [`#76`](https://bitbucket.org/decimalteam/go-node/pull-requests/76) Preprod
-  [`#75`](https://bitbucket.org/decimalteam/go-node/pull-requests/75) feature/23003 fix: counting MissedBlockCounter
-  [`#73`](https://bitbucket.org/decimalteam/go-node/pull-requests/73) hotfix/22943 fix: important optimizations in modules validator, gov and swap.
-  [`#71`](https://bitbucket.org/decimalteam/go-node/pull-requests/71) feature/22158 fix: download upgrades and check hashes for new files
-  [`#55`](https://bitbucket.org/decimalteam/go-node/pull-requests/55) hotfix/14831 fix: sync (forgotten) error wrapping from staging to master
-  [`#58`](https://bitbucket.org/decimalteam/go-node/pull-requests/58) feature/8812 fix: add missed aliases
-  [`#74`](https://bitbucket.org/decimalteam/go-node/pull-requests/74) feature/22265 feat: extended log info for failed tx with code 4
-  [`#70`](https://bitbucket.org/decimalteam/go-node/pull-requests/70) feature/21666 fix: prevent incorrect detection of grace period
-  [`#52`](https://bitbucket.org/decimalteam/go-node/pull-requests/52) removed pruning option of app.conf overwriting and extended grace period to 182 days
-  [`#47`](https://bitbucket.org/decimalteam/go-node/pull-requests/47) Deploy Oracl Linux8
-  [`3394591`](https://bitbucket.org/decimalteam/go-node/commits/33945917546416649e0bdf033a82f65452e8bea9) hotfix/22943-master fix: important optimizations in modules validator, gov and swap.
-  [`e25c1e7`](https://bitbucket.org/decimalteam/go-node/commits/e25c1e7d2bb503751bb86a2e112405fb1d846821) disable centos 8 build
-  [`81bac1f`](https://bitbucket.org/decimalteam/go-node/commits/81bac1f1ae85ca9fc4c7e70d80c0b98798e572e2) up version to 1.2.12

#### [v1.2.10](https://bitbucket.org/decimalteam/go-node/compare/v1.2.10..v1.2.9) - 2022-02-16

-  [`0feeba8`](https://bitbucket.org/decimalteam/go-node/commits/0feeba844916445a17c8a1bcfc9a177e9a8e0663) up version to 1.2.10
-  [`91cac29`](https://bitbucket.org/decimalteam/go-node/commits/91cac290ad38f8f0ab94999493872133d281b036) updated one hour of grace period duration

#### [v1.2.9](https://bitbucket.org/decimalteam/go-node/compare/v1.2.9..v1.2.8) - 2022-02-16

-  [`#30`](https://bitbucket.org/decimalteam/go-node/pull-requests/30) Updating deploy Centos 8
-  [`#29`](https://bitbucket.org/decimalteam/go-node/pull-requests/29) Preprod
-  [`bf21102`](https://bitbucket.org/decimalteam/go-node/commits/bf2110299f155a14b34e00626413b5d83b84eb66) up version to 1.2.9
-  [`150e8e4`](https://bitbucket.org/decimalteam/go-node/commits/150e8e46a0276d9cec71c87774468532ac60472d) Revert "up version to 1.2.9"
-  [`cf6fd2d`](https://bitbucket.org/decimalteam/go-node/commits/cf6fd2ddcef54f0c12d310ec97f12c27096a5afa) Merged master into preprod

#### [v1.2.8](https://bitbucket.org/decimalteam/go-node/compare/v1.2.8..v1.2.7) - 2022-02-04

-  [`1c68c15`](https://bitbucket.org/decimalteam/go-node/commits/1c68c156628f10b7c821699dc304cf62740a978b) up version to 1.2.8
-  [`b0826f6`](https://bitbucket.org/decimalteam/go-node/commits/b0826f6a32ce6006c49f6e713818dfdf0f297313) extended grace period

#### [v1.2.7](https://bitbucket.org/decimalteam/go-node/compare/v1.2.7..v1.2.6) - 2022-01-26

-  [`5afd8c8`](https://bitbucket.org/decimalteam/go-node/commits/5afd8c8b26b86126217c908dd9a3912e6e35cb8f) added build on ubuntu 22.04
-  [`d623585`](https://bitbucket.org/decimalteam/go-node/commits/d623585291591e04a4576807f4d31a2149630a36) update version
-  [`81a363b`](https://bitbucket.org/decimalteam/go-node/commits/81a363b605038abb43ed58803fef83b1a4304f75) extended grace period and fixed missing blocks count

#### [v1.2.6](https://bitbucket.org/decimalteam/go-node/compare/v1.2.6..v1.2.5) - 2022-01-26

-  [`66491a3`](https://bitbucket.org/decimalteam/go-node/commits/66491a394d5ec57d42bd6d712b9fc92110150de9) up cosmos-sdk 0.39.3
-  [`5e1db43`](https://bitbucket.org/decimalteam/go-node/commits/5e1db4371c2deb062ca5862bbd43e8a650bde206) Update nft reserver
-  [`cb72159`](https://bitbucket.org/decimalteam/go-node/commits/cb72159a9b79683c77bc0edac059814acf252219) append sort(unbond_nft), change(missed_blocks), sync(delegated_coin)

#### [v1.2.5](https://bitbucket.org/decimalteam/go-node/compare/v1.2.5..v1.2.4) - 2021-12-14

-  [`e543f46`](https://bitbucket.org/decimalteam/go-node/commits/e543f464ba3788738efa5f7bb609f62fe610df3e) update
-  [`b7e00be`](https://bitbucket.org/decimalteam/go-node/commits/b7e00be874bcecd9685bb3e86fc99eb7c54c609e) from master
-  [`5df6c27`](https://bitbucket.org/decimalteam/go-node/commits/5df6c27c153c2032e020074519b8b184401025a2) fix nft

#### [v1.2.4](https://bitbucket.org/decimalteam/go-node/compare/v1.2.4..v1.2.3) - 2021-12-14

-  [`6497129`](https://bitbucket.org/decimalteam/go-node/commits/6497129b3cb247ee4c35200747ae4a1f5b5b5a21) changes to master
-  [`fcbfe1d`](https://bitbucket.org/decimalteam/go-node/commits/fcbfe1df7c4261095592d01db0583f9ac5ce1d20) from master
-  [`deea294`](https://bitbucket.org/decimalteam/go-node/commits/deea294a134d3c65c634a9b2a5e340f45c063b37) improved readme

#### [v1.2.3](https://bitbucket.org/decimalteam/go-node/compare/v1.2.3..v1.2.2) - 2021-12-14

-  [`83d0e24`](https://bitbucket.org/decimalteam/go-node/commits/83d0e249a4b72439c6f3a6bb953a234955136654) update v1.2.1
-  [`cac414d`](https://bitbucket.org/decimalteam/go-node/commits/cac414dec394a2c0cc646890138ff5c596a3e67b) Add rest for gov module
-  [`71e5110`](https://bitbucket.org/decimalteam/go-node/commits/71e5110be36754079ce4c6d660242757a2e4e5ca) UpdateCoin tx, identity for coin, remove validatorCache

#### [v1.2.2](https://bitbucket.org/decimalteam/go-node/compare/v1.2.2..v1.2.1) - 2021-12-14

-  [`c85c254`](https://bitbucket.org/decimalteam/go-node/commits/c85c254158df7dfd8a8a8e5c3cf79344f5c4cd5b) clone from master
-  [`166700e`](https://bitbucket.org/decimalteam/go-node/commits/166700e223a8efec8611e04d5475272a73541972) add check coin exists , fix delegate , check recv in send coin
-  [`48112dd`](https://bitbucket.org/decimalteam/go-node/commits/48112ddf68d43e4e21b7fe9a6e15ec51ff9c2772) from master

#### [v1.2.1](https://bitbucket.org/decimalteam/go-node/compare/v1.2.1..1.1.14) - 2021-10-28

-  [`b769bb5`](https://bitbucket.org/decimalteam/go-node/commits/b769bb5724fe64ec63f072403ea92dfd2976ef75) from master
-  [`2dc990d`](https://bitbucket.org/decimalteam/go-node/commits/2dc990deb15c91434f7e1594c7efbc2e4745c019) -&gt;master
-  [`0c5664c`](https://bitbucket.org/decimalteam/go-node/commits/0c5664c84b63df51cbc66771f27aa7587a29a472) merge

#### [1.1.14](https://bitbucket.org/decimalteam/go-node/compare/1.1.14..v1.0.4) - 2021-07-27

-  [`#4`](https://bitbucket.org/decimalteam/go-node/pull-requests/4) Feature/5835
-  [`34ac97d`](https://bitbucket.org/decimalteam/go-node/commits/34ac97d8426d6bf382811af8baaf6f3e59e76ebf) Revert
-  [`7d95869`](https://bitbucket.org/decimalteam/go-node/commits/7d958691d67890b080231392ba6a955980201694) Revert bulk to beef2f76
-  [`b00dc2e`](https://bitbucket.org/decimalteam/go-node/commits/b00dc2e764a97a20478e0b06baf68cfee5ef3b34) Revert "Revert bulk to beef2f76"

#### [v1.0.4](https://bitbucket.org/decimalteam/go-node/compare/v1.0.4..v1.0.3) - 2020-08-06

-  [`2bd760b`](https://bitbucket.org/decimalteam/go-node/commits/2bd760bbf2066359df817efb088044aa5a44f867) Fix pay rewards after edit-candidate with new reward address

#### [v1.0.3](https://bitbucket.org/decimalteam/go-node/compare/v1.0.3..v1.0.2) - 2020-08-05

-  [`428fcc5`](https://bitbucket.org/decimalteam/go-node/commits/428fcc5a24742c48aee025881ea02ffd702b30f0) [reset]
-  [`6d5cbd9`](https://bitbucket.org/decimalteam/go-node/commits/6d5cbd97632e544ab9446fd3c16d4fa8a7b6590f) Fix gasUsed = gasWanted, count validator
-  [`5c4d53c`](https://bitbucket.org/decimalteam/go-node/commits/5c4d53c4d1b8960d529b0f65f2bfde664dabc548) GasWanted = GasUsed

#### [v1.0.2](https://bitbucket.org/decimalteam/go-node/compare/v1.0.2..v1.0.0) - 2020-08-03

-  [`4e35c61`](https://bitbucket.org/decimalteam/go-node/commits/4e35c61c2e00deacf04693f23519633cdadfce93) Fixed problem with redeem checks and creating coins.
-  [`b788968`](https://bitbucket.org/decimalteam/go-node/commits/b788968793823d04aef31e0c746285329783f6ce) Fixed problem with redeem checks and creating coins.
-  [`29b8ff8`](https://bitbucket.org/decimalteam/go-node/commits/29b8ff8d84a63a6a91354fc6d1a88f4f2aed5aa6) Updated block for update behavior.

#### v1.0.0 - 2020-08-01

-  [`1e21181`](https://bitbucket.org/decimalteam/go-node/commits/1e21181727cb0bc30a29eec2ee28c5d8dfe5f201) removed unused utils/crypto package
-  [`3054d95`](https://bitbucket.org/decimalteam/go-node/commits/3054d95db6dbc1e382fb24b243ec209a0c1f66c5) Revert "Add the ability to pay a commission with a custom coin"
-  [`4cabe9d`](https://bitbucket.org/decimalteam/go-node/commits/4cabe9dc42b43d67c4844a08999d0af044db7d50) Add the ability to pay a commission with a custom coin
