import { createMuiTheme } from "@material-ui/core/styles";
import { indigo, red } from "@material-ui/core/colors";

const theme = createMuiTheme({
  palette: {
    primary: indigo,
    error: red
  },
  props: {
    MuiPaper: {
      elevation: 4
    }
  }
});

export default theme;
