# l4go/recode ライブラリ

テキストベースのフォーマットを再帰的にデコードするためのライブラリです。  
主な用途は、外部設定ファイルを使った設定ファイルの再帰的デコードです。

## ライブラリーの目的

JSON、YAML、TOMLなどのテキストベースの汎用フォーマットが、設定ファイルに良く使われています。  
フォーマットに汎用性があるので、
後からのパラメータ追加も容易で、機能追加とともに設定ファイルも大きくなりがちです。
機能が増えることはいいのですが、設定ファイルが大きくなると、利用者にとっては見通しが悪くなります。  
このような場合、設定の一部を外部設定ファイルとして切り出せれば、見通しの悪さを軽減できます。
ですが、汎用デコーダは、設定ファイル専用ではないので、外部設定ファイルを想定していません。

そこで、外部設定ファイルへの対応できるように、
汎用デコーダでの再帰的デコードをサポートするのが、
このライブラリの`RecursiveRebuild()`です。


## `RecursiveRebuild()`の動作

親要素をデコードした後、`recode.RecursiveRebuild()`を実行すると再帰的なデコードが行われます。  
ただ、子孫の要素のデコード方法は、`RebuildByType()`メソッドとして定義しておく必要があります。
また、`RebuildByType()`メソッドは、
引数の型を`recode.RecursiveRebuild()`の2番目の引数と一致させておく必要もあります。

また、`RebuildByType()`メソッドがある同じ型を、上位と子孫の要素に同時指定するのは避けてください。
このようにすると、同じ`RebuildByType()`メソッドが呼ばれるので、無限ループします。

`RecursiveRebuild()`では、以下のような処理を行います。

1. 引数が同じ`RebuildByType()`メソッドが存在するフィールド値を探す。
    - 該当するフィールドがなければ終了する。
2. `RebuildByType()`メソッドを呼び出す。
3. `RebuildByType()`メソッドが呼び出されたフィールド値に、1.からの処理を再帰的に実行する。

## `RecursiveRebuild()`と`RebuildByType()`の引数の型

`RecursiveRebuild()`の定義は以下のようになっています。

```go
func RecursiveRebuild[T any](v any, param T) error
```

`RebuildByType()`を定義する、`AnyRebuilder`インタフェイスは以下のようになっています。

```go
type AnyRebuilder[T any] interface {
	RebuildByType(param T) error
}
```

2つの定義の両方にある`param T`の引数に、Goのジェネリックス(generics)の機能を使っていますので、
型を自由に変えられます。
ただし、`param T`部分の引数の型を、 `RecursiveRebuild()`関数と`RebuildByType()`メソッドで、そろえる必要があります。
引数の型が一致してないと、`RebuildByType()`が呼び出されなくなります。

## サンプル概要

Example形式の [サンプルプログラム](../ex_test.go)は、以下の処理を行っています。

- [testfsディレクトリ](../testfs)を`os.DirFS()`でfs.FS化
    - 例えば、embed.FSでも同じことが出来ます。
- fs.FS経由でファイルを読み取って、再帰的にデコード

このデコードでは以下の処理が再帰的に行われます。

1. [textfile.json](../testfs/testfile.json)をTextTextFile型にデコード
2. [text.json](../testfs/test.json)をTextTextFileのTextFileフィールドにデコード
3. [text.txt](../testfs/test.txt)の内容をTextFileのTextフィールドに設定
