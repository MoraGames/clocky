<ins>**v0.4.1 - It's not the v2, Balancing Update 1**</ins>

**Rebalance number of less/more impacting effects:**
- ~~Remove Add1, Sub1 and Mul-1 effects~~
- ~~Rebalance presence of the other effects~~
- ~~Remove Reigning Podium effect~~
- ~~Reduce to +1 the bonus of Reigning Leader effect~~

**Rebalance the distribution of Events grouped by Sets:**
- ~~Remove Equals set~~
- ~~Implements Equal Twins set (aa:bb)~~
- ~~Remove a!=b constraint in Repeat and Mirror sets~~
- ~~Implements Half set (n:n/2)~~

**Fix Championship visualization bugs:**
- ~~Edit the command /stats~~
- ~~Edit the command /ranking and related GetRanking function~~
- ~~Implements reset message in the group chat~~

**Hints and Activity Streak buff:**
- ~~Reduce the requirments to be considered a "participant user"~~
- ~~Link hints to the partecipation (not the activity)~~
- ~~Produce the hints message with 3 different sets (that must have at least a total of 20 events)~~
- ~~Add and rebalance the bonus given by Activity Streak (one for each 7 days: 7, 14, 21, 28, ...)~~

**Other:**
- ~~Resolve cronjob 2-weekly WithStartDateTimePast() execute instantly~~
- Implements summary message with some stats on daily time base
- Implements summary message with some stats on championship time base

<ins>**v0.4.2 - It's not the v2, Balancing Update 2**</ins>

**Rebalance number of less/more impacting effects:**
- Slightly reduce amount of Add3, Add4, Sub3, Sub4
- Reduce amount of Mul2, Mul3, Mul4, Mul5, Mul6, Mul-2, Mul-3, Mul-4 and Mul-5 effects

**Other:**
- ~~Rewrite the entire message entity parser~~
- ~~Rewrite the handling of commands, their descriptions and their responses~~