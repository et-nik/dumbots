# [DumBots](https://github.com/et-nik/dumbots)

This is simple bots for Half-Life 1 game server.
They are suitable for plugin developers to test plugin functionality. 
The bots can execute simple commands, move to a designated point, and overcome obstacles.

Bots written by with [Metamod-Go](https://github.com/et-nik/metamod-go) library.

## Bot commands

* `dbt_list` — list of bots
* `dbt_add` — add bot
* `dbt_add_mock` — add a bot that won't do anything
* `dbt_move <id> <x> <y> <z>` — teleport bot to the specified coordinates
* `dbt_goto <id> <x> <y> <z>` — add bot task to go to the specified coordinates
* `dbt_angles <id> <pitch> <yaw>` — set bot angles
* `dbt_rm <id>` — remove bot
* `dbt_rm_all` — remove all bots
* `dbt_press <id> <key>` — press key (attack, attack2, jump, duck, use, reload)
* `dbt_room_graph <id>` — show room graph for bot
* `dbt_hide_room_graph <id>` — hide room graph for bot
* `dbt_say <id>` — show room for bot

## CVars

* `dbt_max_bots` — maximum number of bots