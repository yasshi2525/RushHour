import React, {
  Suspense,
  PropsWithChildren,
  lazy,
  useEffect,
  useMemo
} from "react";
import ThemeProvider from "@material-ui/styles/ThemeProvider";
import { SnackbarProvider, SnackbarProviderProps } from "notistack";
import LoadingCircle from "common/utils/loading";
import { OperationProvider } from "common/utils/operation";
import { LoadingProvider } from "common/loading";
import { AdminPageProvider } from "common/admin";
import theme from "./theme";
import LoadingProgress from "./Loading";
import Maintenance from "./Maintenance";

const Application = lazy(() => import("./Application"));

const snackOpts: SnackbarProviderProps = {
  maxSnack: 5,
  anchorOrigin: {
    vertical: "bottom",
    horizontal: "right"
  }
};

/**
 * `localStorage["jwt"]` からユーザ情報の取得を試みて、コンポーネントを描画する
 */
export default (props: PropsWithChildren<{ admin?: boolean }>) => {
  useEffect(() => {
    console.info("after RootElement");
  }, []);
  return useMemo(
    () => (
      <ThemeProvider theme={theme}>
        <SnackbarProvider {...snackOpts}>
          <LoadingProvider>
            <LoadingProgress />
            <AdminPageProvider admin={props.admin}>
              <OperationProvider maintenance={<Maintenance />}>
                <Suspense fallback={<LoadingCircle />}>
                  <Application />
                </Suspense>
              </OperationProvider>
            </AdminPageProvider>
          </LoadingProvider>
        </SnackbarProvider>
      </ThemeProvider>
    ),
    []
  );
};
