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

拡張手順
=========

ここではRushHourをエンハンスするときの追加手順を記述します。

デフォルトパラメータの追加
--------------------------

定数として用いる値を追加する手順を示します。

プロパティファイルの編集
^^^^^^^^^^^^^^^^^^^^^^^^

プロパティを追加します。単体テストでも使うため、テスト用プロパティも追加します。

`main/resources/template_config.properties`
`test/resources/template_config.properties`

例::

    rushhour.game.default.platform.capacity=20

プロパティ管理クラスの編集
^^^^^^^^^^^^^^^^^^^^^^^^^^

プロパティ管理クラスに定数を追加します。

`class net.rushhourgame.RushHourProperties`

例::

    public static final String GAME_DEF_GATE_NUM = "rushhour.game.default.platform.capacity";

呼び出し方法
^^^^^^^^^^^^

RushHourPropertiesのインスタンス prop とすると

例:: 
    
    Integer.parseInt(prop.get(RushHourProperties.GAME_DEF_GATE_NUM));

Controllerの追加
-----------------

`class net.rushhourgame.controller.AbstractController` を継承します。

単体テスト
^^^^^^^^^^^

`class net.rushhourgame.controller.ControllerFactory` に
`create追加するコントローラ名` メソッドを追加します。

`class net.rushhourgame.controller.XXXControllerTest` を追加し、
`class net.rushhourgame.controller.AbstractControllerTest` を継承します。

XXXControllerTestに::

    protected final static XXXController XCON = ControllerFactory.createXXXController();
    