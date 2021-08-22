export const getToken = () => {
  return sessionStorage.getItem("redux") &&
    JSON.parse(sessionStorage.getItem("redux")).storage.token
    ? JSON.parse(sessionStorage.getItem("redux")).storage.token
    : false;
};

export const parsingToken = (token: string) => {
  const tokenSplit = token.split(".")[1];
  const tokenDecode = decodeURIComponent(
    atob(tokenSplit)
      .split("")
      .map(function(c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );
  return JSON.parse(tokenDecode);
};

export const getUserInfo = () => {
  const token = getToken();
  const admin = parsingToken(token).userType === "USER_ADMIN" ? true : false;
  const userPermission = parsingToken(token).permission;
  return { admin, permission: userPermission };
};
