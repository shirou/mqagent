


Actions
----------------

- metric
- spec
- check
- command



TODO
-------

- duplicated actionid between subscribes
- stop sending metrics

  - required when subscribed roles are reloaded

- include username to topic root



example
-----------------

<hostid> := some random 8byte string.

shirou@github/mqagent/abcdefgh/load/avg5 -> abcdefgh.load.avg5
shirou@github/mqagent/abcdefgh/system/ -> abcdefgh.system.mem
shirou@github/mqagent/abcdefgh/hostinfo/os -> abcdefgh.hostinfo.linux



