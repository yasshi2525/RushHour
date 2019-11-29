import React, { createContext } from "react";
import { ComponentProperty } from "interfaces/component";

const AdminPageContext = createContext(false);

interface ProviderProperty extends ComponentProperty {
  admin: boolean;
}
export const AdminPageProvider = (props: ProviderProperty) => {
  return (
    <AdminPageContext.Provider value={props.admin}>
      {props.children}
    </AdminPageContext.Provider>
  );
};

export default AdminPageContext;
