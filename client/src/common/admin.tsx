import React, { FC, createContext, useMemo, useEffect } from "react";

const AdminPageContext = createContext(false);
AdminPageContext.displayName = "AdminPageContext";

export const AdminPageProvider: FC<{ admin?: boolean }> = props => {
  useEffect(() => {
    console.info("after AdminPageProvider");
  }, []);
  return useMemo(
    () => (
      <AdminPageContext.Provider value={props.admin === true}>
        {props.children}
      </AdminPageContext.Provider>
    ),
    []
  );
};

export default AdminPageContext;
