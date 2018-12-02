# github-developer-rate

A program made with Golang to help us discover how much reviews a user did in some repositories.

Install

go build rate.go pull_request.go pull_review.go

Use

./rate -owner=accountname -repo=repositoryname,repositoryname2 -dateStart=2018-10-01 -dateEnd=2018-11-20 -githubUser=loginname

Example

./rate -owner=guilhermesiani -repo=workspace,interactivehistory,virtualpet, -dateStart=2018-10-01 -dateEnd=2018-11-20 -githubUser=chewbacca

Steps on execute

1) Put your github user
2) Put your github password

Output should be something like:

```
Calculating 

workspace
..........
interactivehistory
....
virtualpet
.......

User guilhermesianilinx did 6 reviews about 16 pull requests (36%) on interval between 2018-11-10T00:00:00Z and 2018-11-30T23:59:59Z
Researchs ends for chaordic [newoffers-inventory chubaca oms-stock-notifier oms-stock-threshold-importer falcon freight-importer]
```
