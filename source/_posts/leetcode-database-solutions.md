title: LeetCode 数据库题目解答
date: 2016-04-06 16:02:00
categories:
- LeetCode
tags:
- LeetCode
- 数据库
- SQL
toc: true

---

前两天重刷了《SQL必知必会》，昨天想到了 LeetCode，于是去刷了几道数据库的题目，开了不少脑洞。

<!-- more -->

今天把答案整理一下。喔，题库在[这里](https://leetcode.com/problemset/database/)。

# 175. Combine Two Tables

<https://leetcode.com/problems/combine-two-tables/>

样例中有些人的 `PersonId` 无法在 `Address` 表中找到，所以使用 `LEFT JOIN`.

```SQL
SELECT FirstName, LastName, City, State
FROM Person
LEFT JOIN Address
ON Person.PersonId = Address.PersonId;
```

# 176. Second Highest Salary

<https://leetcode.com/problems/second-highest-salary/>

`UNION` 查询，在结果的最后添加一个 `NULL`, 若不存在第二高的薪水则会选择 `NULL`.

```SQL
SELECT Salary FROM Employee
UNION
SELECT NULL
ORDER BY Salary DESC
LIMIT 1, 1;
```

# 177. Nth Highest Salary

<https://leetcode.com/problems/nth-highest-salary/>

这个不知道为什么不可以用 `LIMIT 1, N-1`，所以用了 `IF` 函数。

```SQL
CREATE FUNCTION getNthHighestSalary(N INT) RETURNS INT
BEGIN
  RETURN (
      # Write your MySQL query statement below
      SELECT IF(COUNT(*) >= N, MIN(rank.Salary), NULL)
      FROM (
        SELECT DISTINCT Salary
        FROM Employee
        ORDER BY Salary DESC
        LIMIT N
    ) AS rank
  );
END
```

# 178. Rank Scores

<https://leetcode.com/problems/rank-scores/>

```SQL
SELECT Score, (
    SELECT COUNT(DISTINCT Score)
    FROM Scores AS c
    WHERE o.Score <= c.Score      # 统计比已选分数小的分数个数
) AS Rank
FROM Scores AS o
ORDER BY Score DESC;
```

# 180. Consecutive Numbers

<https://leetcode.com/problems/consecutive-numbers/>

暴力查询:joy:

```SQL
SELECT DISTINCT l1.Num AS ConsecutiveNums
FROM Logs AS l1, Logs AS l2, Logs AS l3
WHERE l1.Id+1 = l2.Id AND l2.Id+1 = l3.Id
  AND l1.Num = l2.Num AND l2.Num = l3.Num;

```

# 181. Employees Earning More Than Their Managers

<https://leetcode.com/problems/employees-earning-more-than-their-managers/>

选择雇员，根据 `ManagerId` 找到雇员上司的薪水，然后进行比较即可。

```SQL
SELECT Name
FROM Employee
WHERE Salary > (
          SELECT Salary
          FROM Employee AS e
          WHERE e.id = Employee.ManagerId
        );
```

# 182. Duplicate Emails

<https://leetcode.com/problems/duplicate-emails/>

按 `Email` 字段进行分类，使用 `HAVING` 筛选出相同 `Email` 数量大于 1 的项。

```SQL
SELECT Email FROM Person
GROUP BY Email
HAVING COUNT(Email)>1;
```

# 183. Customers Who Never Order

<https://leetcode.com/problems/customers-who-never-order/>

这个也是直接查询…

```SQL
SELECT c.Name AS Customers
FROM Customers AS c
WHERE (SELECT COUNT(*)
       FROM Orders
       WHERE c.id = Orders.CustomerId) = 0;
```

# 184. Department Highest Salary

<https://leetcode.com/problems/department-highest-salary/>

基本上就是直接查询，注意 `WHERE` 语句中判别条件的位置，否则有可能 TLE:joy:

```SQL
SELECT d.Name AS Department,
    e.Name AS Employee,
    e.Salary
FROM Employee AS e, Department AS d
WHERE e.DepartmentId = d.Id
  AND e.Salary = (SELECT MAX(e2.Salary)
                FROM Employee AS e2
                WHERE e.DepartmentId = e2.DepartmentId);
```

# 185. Department Top Three Salaries

<https://leetcode.com/problems/department-top-three-salaries/>

输出每个部门薪资最高的三个人。这个题里有个坑，如果两个人薪资相同，那么这两个人并列，都要输出。并且如果四个人的薪资为 3 2 2 1， 薪资为 1 的那个人排第 3 :joy:

```SQL
SELECT d.Name AS Department,
    e.Name AS Employee,
    e.Salary
FROM Employee AS e, Department AS d
WHERE e.DepartmentId = d.Id
  AND (SELECT COUNT(DISTINCT e2.Salary)             # 排序时允许并列
            FROM Employee AS e2
            WHERE e.DepartmentId = e2.DepartmentId
              AND e.Salary < e2.Salary) < 3         # 比该雇员工资高的人少于三个
ORDER BY Department, Salary DESC;
```

# 196. Delete Duplicate Emails

<https://leetcode.com/problems/delete-duplicate-emails/>

MySQL 不允许在删除时依据待删除的表进行筛选 (You can't specify target table'Person' for update in FROM clause), 所以要绕一下。

```SQL
# 错的！！
DELETE FROM Person
WHERE Id IN (SELECT p2.Id
             FROM Person AS p1, Person AS p2
             WHERE p1.Email = p2.Email
               AND p1.Id < p2.Id
            );

DELETE FROM Person
WHERE Id IN (SELECT * FROM(                         # 绕一下，先挑出所有满足要求的 ID 构成一个表，再从这个表中选 Id 进行删除
                    SELECT p2.Id
                    FROM Person AS p1, Person AS p2
                    WHERE p1.Email = p2.Email
                      AND p1.Id < p2.Id) AS temp
                      );
```

# 197. Rising Temperature

<https://leetcode.com/problems/rising-temperature/>

主要考 `MySQL` 的日期操作函数。

```SQL
SELECT w1.Id AS Id
FROM Weather AS w1, Weather AS w2
WHERE datediff(w1.Date, w2.Date) = 1
  AND w1.Temperature > w2.Temperature;
```

# 262. Trips and Users

<https://leetcode.com/problems/trips-and-users/>

太乱了，没做:disappointed_relieved:

# 后记

昨天花了半天写完这些题，写到最后都不知道自己在写什么了:joy:不过还是掌握了不少的 SQL 查询技巧，比如 `UNION SELECT NULL` 等等。
