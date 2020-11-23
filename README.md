# diceroller

A tool for generating fair or biades dice-rolls for an arbitrary sided dice.
It supports command line and can be started as a [rest service](#api).

## Command line

Throw a fair 6-sided dice one time
```
.\diceroll.exe 
```

Throw a fair 6-sided dice 10 times
```
.\diceroll.exe -rolls 10
```

Throw a biased 3-sided dice 6 times
```
.\diceroll.exe -rolls 6 -odds "1,1,2"
.\diceroll.exe -rolls 6 -probs "1/4,1/4,2/4"
```


## API

Compile and run with -srv flag.
For example in windows

```
diceroller -srv
```

Optionally a port can be supplied

```
diceroller -srv -port 8080
```


Default port is 10000

Request format
```
http://localhost:10000/:rolls
http://localhost:10000/:rolls/odds/:data
http://localhost:10000/:rolls/probs/:data
```

`:rolls` is the integer that definies the amount of dice rolls requested

Example of requests:

```
http://localhost:10000/1
http://localhost:10000/100
```


Optionally the number of dice-sides and their odds or probablitites can be set up for the dice rolls. `:data` is a comma separated list of the sides and their respective odds or propabilites.

Example a three-sided dice where side 1 & 2 have a probability of 0.25 and side 3 have a probability of 0.5.
```
http://localhost:10000/10/odds/1,1,2
http://localhost:10000/10/probs/0.25,0.25,0.5

```






