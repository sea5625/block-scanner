import { call, put, takeLatest } from "redux-saga/effects";
import service from "../helper/service";
import { actions as storageActions } from "modules/storage";
import { getToken, parsingToken } from "utils/token";
import { history } from "store";
import { USER_TYPE_THIRD_PARTY } from "utils/const";
import Popup from "components/common/popup/Popup";

//authActionTypes
const LOGIN = "LOGIN";
const LOGIN_SUCCESS = "LOGIN_SUCCESS";
const LOGIN_FAILURE = "LOGIN_FAILURE";

const LOGOUT = "LOGOUT";
export const LOGOUT_SUCCESS = "LOGOUT_SUCCESS";

//authPayloadTypes
interface AuthState {}

interface LoginPayload {
    id: string;
    password: string;
}

interface LoginData {
    token: string;
    passwordStatus: number;
}

//authActions
export const actions = {
    login: (payload: LoginPayload) => ({
        type: LOGIN,
        payload
    }),
    loginSuccess: () => ({
        type: LOGIN_SUCCESS
    }),
    loginFailure: (errorMessage: string) => ({
        type: LOGIN_FAILURE,
        errorMessage
    }),
    logout: () => ({
        type: LOGOUT
    })
};

//authReducer
export function authReducer(
    state: AuthState = {
        errorMessage: ""
    },
    action
): AuthState {
    switch (action.type) {
        case LOGIN_FAILURE:
            const { errorMessage } = action;
            return {
                ...state,
                errorMessage
            };
        default:
            return state;
    }
}

//authAPI
const api = {
    login: async (payload: LoginPayload) => {
        return await service.post("/auth/login", payload);
    }
};

//authSaga
function* loginFunc(action) {
    try {
        const payload: LoginPayload = action.payload;
        const res: LoginData = yield call(api.login, payload);
        if (res) {
            const { token, passwordStatus } = res;
            const tokenParse = parsingToken(token);
            if (tokenParse.userType === USER_TYPE_THIRD_PARTY) {
            } else {
                yield put(storageActions.setToken(token));
                yield put({
                    type: LOGIN_SUCCESS
                });
                if (passwordStatus === 3002) {
                    const id = tokenParse.user.toString();
                    Popup.changePasswordPopup({
                        id,
                        className: "change-password"
                    });
                } else {
                    history.push("/");
                }
            }
        }
    } catch (e) {
        const { internalMessage: errorMessage } = e.data.errors[0];
        yield put({
            type: LOGIN_FAILURE,
            errorMessage
        });
    }
}

function* logoutFunc() {
    try {
        yield put(storageActions.removeToken());
        if (!getToken()) {
            yield put({ type: LOGOUT_SUCCESS });
            history.push("/");
        }
    } catch (e) {
        console.log(e);
    }
}

export function* authSaga() {
    yield takeLatest(LOGIN, loginFunc);
    yield takeLatest(LOGOUT, logoutFunc);
}
