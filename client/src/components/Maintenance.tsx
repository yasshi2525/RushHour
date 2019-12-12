import React, { useMemo, Fragment, useContext } from "react";
import AdminPageContext from "common/admin";
import LogOut from "./LogOut";

const AdminPage = () => {
  return (
    <Fragment>
      <LogOut />
      [TODO]管理者操作
    </Fragment>
  );
};

const TopPage = () => (
  <div>メンテナンス中です。時間をあけて再アクセスしてください。</div>
);

export default () => {
  const isAdminPage = useContext(AdminPageContext);

  return useMemo(() => (isAdminPage ? <AdminPage /> : <TopPage />), [
    isAdminPage
  ]);
};
