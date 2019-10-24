.. MIT License

    Copyright (c) 2017 yasshi2525

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

データモデル
============

データモデルを示します。

設計方針
--------

* 主キーは人工キーではなくナチュラルキーで定義する
* 複合主キーになる場合、サロゲートキーを定義する
* 外部キーはサロゲートキーに対して設定する
* `nullable` なフィールドは定義しない

論理データモデル
----------------

.. note::

    Ver 0.0では人の経路情報をデータベースに永続化したが、
    `nullable` なフィールドが増え、条件分岐が複雑になってしまった。
    また、人ごとに経路情報を持つため、パフォーマンスが出なかった。
    そこで Ver 1.0 は性能面の課題を解決するため、経路情報は
    ゲームマスタが管理する仕様とした。

.. note::

    人用移動ステップはデータベースに保存しなくとも、実行時に生成可能。
    しかし、永続化する情報によって決まる値なので、永続化対象にした。
