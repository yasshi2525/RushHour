import * as React from "react";
import { connect } from "react-redux";
import { ThemeProvider } from "@material-ui/styles";
import { GameBoardProperty } from "../common/interfaces";
import ResizeHandler from "../common/handlers/window"
import { RushHourStatus } from "../state";
import { initPIXI, players, fetchMap } from "../actions";
import RushHourTheme from ".";
import Canvas from "./Canvas";
import ToolBar from "./Toolbar";
import GameBar from "./GameBar";

// ゲーム画面のルートコンポーネント
export class GameBoard extends React.Component<GameBoardProperty, RushHourStatus> {
    resize: ResizeHandler;
    
    constructor(props: GameBoardProperty) {
        super(props);

        this.resize = new ResizeHandler(props.game.model, this.props.dispatch);
    }
    
    componentDidMount() {
        this.props.dispatch(initPIXI.request(this.props.game));
    }

    componentDidUpdate() {
        if (this.props.isPIXILoaded) {
            if (!this.props.isPlayerFetched) {
                this.props.dispatch(players.request({ model: this.props.game.model, dispatch: this.props.dispatch }));
            } else {
                this.props.dispatch(fetchMap.request({ model: this.props.game.model, dispatch: this.props.dispatch }));
            }
        }
    }

    render () {
        return (
            <ThemeProvider theme={RushHourTheme}>
                { this.props.readOnly ? 
                    <GameBar readOnly={true} />
                : <GameBar readOnly={false} displayName={this.props.displayName} image={this.props.image} /> }
                { this.props.isPIXILoaded ? 
                    <>
                        <Canvas readOnly={this.props.readOnly} model={this.props.game.model} />
                        <ToolBar readOnly={this.props.readOnly} model={this.props.game.model} />
                    </>
                : "ロード中" }
            </ThemeProvider>
        );
    }
}

function mapStateToProps(state: RushHourStatus) {
    return { 
        isPIXILoaded: state.isPIXILoaded,
        isPlayerFetched: state.isPlayerFetched,
    };
}

export default connect(mapStateToProps)(GameBoard);