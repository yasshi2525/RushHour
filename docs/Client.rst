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

クライアント
============

線路敷設の流れ
--------------

線路の敷設を開始し、延伸する

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;

        ユーザ => gameview [label = "onDragEnd()", leftnote = "画面をクリック"] {
            gameview ->> JSF [label = "registerClickPos()", note = "ここでダイアログを開いても\nclose時イベント処理できない"];
        };
        JSF => GameViewBean [label = "registerClickPos()", note = "クリック座標を保存"];
        JSF ->> gameview [label = "fireClickMenu()"];
        gameview ->> GameViewBean [label = "openClickMenu()", note = "close時イベント処理するため、\ncommandButtonから発火"];
        GameViewBean -> clickmenu [label = "openDiaglog()", leftnote = "ダイアログ表示"];
    }

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;
        ユーザ => clickmenu [label = "線路敷設選択"] {
            clickmenu ->> ClickMenuBean [label = "createRail()"];
        }
        ClickMenuBean ->> gameview [label = "closeDialog()", note = "JSはclickmenu上で動かないので\ngameviewに戻す"];
        gameview ->> GameViewBean [label = "handleReturn()"];
        GameViewBean ->> gameview [label = "startExtendingMode(x, y)", leftnote = "線路ポイントを描画"];
    }

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;
        ユーザ => gameview [label = "onDragEnd()", leftnote = "画面をクリック"] {
            gameview ->> JSF [label = "registerClickPos()"];
        };
        JSF => GameViewBean [label = "registerClickPos()", leftnote = "共通処理",  note = "クリック座標を保存"];
        JSF ->> gameview [label = "fireClickMenu()"];
        gameview ->> GameViewBean [label = "extendRail()"];
        GameViewBean ->> gameview [label = "nextExtendingMode(x, y)", leftnote = "描画"];
    }

線路を撤去する。

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;
        ユーザ => gameview [label = "onDragEnd()", leftnote = "線路エッジをクリック"] {
            gameview ->> JSF [label = "registerEdgeId()"];
        };
        JSF => GameViewBean [label = "registerEdgeId", note = "線路エッジのIDを保存"];
        JSF ->> gameview [label = "fireClickMenu()"];
        gameview ->> GameViewBean [label = "openClickMenu()"];
        GameViewBean -> clickmenu [label = "openDiaglog()", leftnote = "ダイアログ表示"];
    }

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;
        ユーザ => clickmenu [label = "線路撤去選択"] {
            clickmenu ->> ClickMenuBean [label = "removeRail()"];
        }
            ClickMenuBean ->> gameview [label = "closeDialog()"];
            gameview ->> GameViewBean [label = "handleReturn()"];
            GameViewBean ->> JSF [label = "PF('confirmDialog').show()", leftnote = "確認ダイアログ表示"];
    }

.. seqdiag::

    seqdiag {
        ユーザ; gameview; clickmenu; JSF; GameViewBean; ClickMenuBean;
        ユーザ => JSF [leftnote = "「削除」をクリック"] {
            JSF ->> GameViewBean [label = "removeRail()"];
        }
        GameViewBean ->> JSF [label = "PF('confirmDialog').hide()"];
        JSF -> gameview [label = "handleCompleteRemoving()"];
    }