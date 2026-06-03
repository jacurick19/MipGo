# MipGo

The World's Fastest Mixed Integer Program Modeling Library for Go

MipGo is a library for modeling and solving Mixed-Integer Lienar Programs (MIPs). It is inspired by the great work of the Coin-OR foundation's Python MIP. Along these lines, some of the main features are:

* More modeling components than any other Go modeling library
  - Cut generators and lazy constraints: work with strong formulations with a large number of constraints by generating only the required inequalities during the branch and   cut search;
  - Solution pool: query the elite set of solutions found during the search;
  - MIPStart: use a problem dependent heuristic to generate initial feasible solutions for the MIP search;
  - Automatic Big-M reformulation of SOS1 constraints.
* Blazing fast
  - MipGo interfaces directly with the native executable for each supported solver. The pure Go modeling components provide best in class speed.
* Multi Solver
  - MipGo supports many mainstream commercial and open source solvers including Gurobi, HiGHS, SCIP, and CBC

## Examples

All of the same examples from Python MIP were implemented in MipGo using the same formulations at https://github.com/jacurick19/MipGo/tree/main/examples.
For details, please see Many Python-MIP examples are documented at https://docs.python-mip.com/en/latest/examples.html.

## Benchmarks

The n-queens problem was modeled with several popular modeling libraries. Here are the results:

### n = 1,000

| Library | Modeling Time (s) |
| :--- | :---: |
| **PuLP** | 5.389 |
| **python-mip** | 4.282 |
| **gurobipy** | 3.187 |
| **Julia JuMP** | 0.770 |
| **nextmv-io/go-mip** | 0.132 |
| **MipGo** | 0.078 |

---

### n = 10,000

| Library | Modeling Time (s) |
| :--- | :---: |
| **Julia JuMP** | 81.853 |
| **nextmv-io/go-mip** | 18.473 |
| **MipGo** | 8.462 |

## FAQ

### Why Go?

I like Go. like mathematical programming, but the tooling for doing it in Go wasn't great. Python, Java, and C++ are pretty standard to use for this type of work. But now, if you ever need to optimize something when you're working with my favorite language in the world, this project exists.

### Will you add support for $MY_FAVORITE_SOLVER???

Yes, probably at some point. If it's open source I would be happy to. If it's a commercial solver and they won't give me a license it will be difficult for me to test the correctness. 

### Are any other features planned?

I would like to add support for more modeling constructs, especially those primarily seen in closed source solvers. Things like (bi)conditional statements, piecewise linear functions, etc.
