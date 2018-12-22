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

Controller
==========

Controllerは :term:`プレイヤ` または :term:`ゲームマスタ` の
要求を実現するモジュールです。ビジネスロジック層に該当します。

.. seqdiag::

    seqdiag {
        GameMaster => Controller [label = "someAction()"] {
            Controller -> EntityManager [label = "find()"];
            Controller <-- EntityManager [label = "List<Entity>"];
            Controller => Entity [label = "someAction()"];
            Controller => Entity [label = "someAction()"];
            Controller => EntityManager [label = "persist(Entity)"];
        }
    }

引数のバリデーションは、実装量削減のため、Bean Validationを用います。

路線
-----

路線は路線ステップ発車、路線ステップ移動、路線ステップ停車、路線ステップ通過の4つから構成されます。
それぞれ以下のように状態遷移します。

.. list-table::
    :header-rows: 1

    * - 遷移元/遷移先
      - 発車
      - 移動
      - 停車
      - 通過

    * - 発車
      - ×
      - ○
      - ○
      - ○

    * - 移動
      - ×
      - ○
      - ○
      - ○

    * - 通過
      - ×
      - ○
      - ○
      - ○

    * - 停車
      - ○
      - ×
      - ×
      - ×
