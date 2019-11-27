import { createContext, Dispatch } from "react";

const initialState: [boolean, Dispatch<boolean>] = [false, () => {}];
export default createContext(initialState);
