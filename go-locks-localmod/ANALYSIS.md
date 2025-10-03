# Analysis — Questions 1 and 2

## Q1 (remove guard; assume no wakeup/wait race)
Under the stated assumption, the algorithm continues to ensure mutual exclusion and progress. The removed guard primarily protects against scheduling/wakeup races; if such races are excluded by assumption, the remaining logic (test-and-set with queueing and wakeup) is sufficient. In practice, removing the guard can reduce robustness if the assumption does not hold.

## Q2 (LL/SC: replace `unlock` with `flag = flag - 1`)
Incorrect. The lock relies on the invariant `flag ∈ {0,1}`. Decrementing can create values outside this set, so lockers may never observe `flag == 0`, breaking liveness (possible livelock or deadlock) and invalidating safety reasoning.
