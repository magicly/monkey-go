# 目标

实现一门类似 Javascript 的语言， 解释执行

功能包括： 整数、bool、字符串等类型， +-\*/<>等数学运算， 定义变量， 算术运算， 定义函数， 函数是 first class 所以支持高阶函数

用处：

- 通过自己实现一门语言， 了解语言的特性以及运行原理， 对以后学习新语言起到事半功倍的效果。
- 为 Babel/ESlint 等 写插件，实现一些想要的功能等。
- 一些意想不到的用处， 比如可以在微信小程序里做热更新。

# 特色

1. 没太多理论，基本不需要 compiler 的理论知识
2. 完全够从零开始， 不使用任何第三方库
3. 用 JS（或 TS） 实现（很容易迁移成其他语言，为了方便用 JS 讲解）

# 几大部分

## Lexer

很简单

## Parser

Top Down Operator Precedence(or: Pratt Parsing)

**重点！！！**

**一定要比较熟悉递归， 不熟悉的可以复习一下， 找些题练习一下**

## Evaluation

**一定要比较熟悉递归， 不熟悉的可以复习一下， 找些题练习一下**

Tree-Walking Interpreter

## 库函数

Question:

1. NodeJS 和 javascript 啥关系、区别？
2. setTimeout， Date 等属于 Javascript 吗

## bytecode & VM

不讲

# 参考资料

- [Writing An Interpreter In Go](https://interpreterbook.com/) ， 极好的书， 强烈推荐看完， 我讲的内容，基本就是参考此书
- [Stanford : Compilers](https://www.bilibili.com/video/av27845355/?p=1) 暂时看前三节即可， https://lagunita.stanford.edu/courses/Engineering/Compilers/Fall2014/course/ 此处有字幕， 有作业， 可以提交。 比较理论、硬核。
