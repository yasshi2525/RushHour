import React, { Suspense, lazy, useEffect, useContext } from "react";
import ThemeProvider from "@material-ui/styles/ThemeProvider";
import { SnackbarProvider } from "notistack";
import { ComponentProperty } from "interfaces/component";
import LoadingContext, {
  LoadingStatus,
  LoadingCircle,
  LoadingProvider
} from "common/loading";
import { AdminPageProvider } from "common/admin";
import { AuthProvider } from "common/auth";
import theme from "./theme";
import LoadingProgress from "./Loading";

const Operation = lazy(() => import("./Operation"));

interface RootComponentProperty extends ComponentProperty {
  admin?: boolean;
}

const Contents = () => {
  const { update } = useContext(LoadingContext);
  useEffect(() => {
    console.info(`effect Root.Contents ${LoadingStatus.IMPORTED_OPERATION}`);
    update(LoadingStatus.IMPORTED_OPERATION);
  }, []);
  return (
    <Suspense fallback={<LoadingCircle />}>
      <Operation />
    </Suspense>
  );
};

/**
 * `localStorage["jwt"]` からユーザ情報の取得を試みて、コンポーネントを描画する
 */
export default (props: RootComponentProperty) => {
  return (
    <ThemeProvider theme={theme}>
      <SnackbarProvider
        maxSnack={5}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "right"
        }}
      >
        <LoadingProvider>
          <AdminPageProvider admin={props.admin === true}>
            <LoadingProgress />
            <AuthProvider>
              <Contents />
            </AuthProvider>
          </AdminPageProvider>
        </LoadingProvider>
      </SnackbarProvider>
    </ThemeProvider>
  );
};
