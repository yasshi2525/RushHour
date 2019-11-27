import React, { Suspense, lazy, useState, useEffect } from "react";
import ThemeProvider from "@material-ui/styles/ThemeProvider";
import { ComponentProperty } from "interfaces/component";
import LoadingContext, { LoadingStatus, useLoading } from "common/loading";
import AdministratorContext from "common/admin";
import LoginContext, { useAuth } from "common/auth";
import theme from "./theme";
import LoadingProgress, { LoadingCircle } from "./Loading";

const Operation = lazy(() => import("./Operation"));

interface RootComponentProperty extends ComponentProperty {
  admin?: boolean;
}

const Contents = () => {
  const [, update] = useLoading();
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
  const isAdminPage = props.admin === true;
  const [status, update] = useState<LoadingStatus>(
    LoadingStatus.CREATED_ELEMENT
  );
  const handlers = useAuth();

  return (
    <ThemeProvider theme={theme}>
      <LoadingContext.Provider value={[status, update]}>
        <AdministratorContext.Provider value={isAdminPage}>
          <LoadingProgress />
          <LoginContext.Provider value={handlers}>
            <Contents />
          </LoginContext.Provider>
        </AdministratorContext.Provider>
      </LoadingContext.Provider>
    </ThemeProvider>
  );
};
