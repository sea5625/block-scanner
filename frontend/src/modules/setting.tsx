import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";
import { Popup } from "components";

const GET_SETTING = "GET_SETTING";
const GET_SETTING_SUCCESS = "GET_SETTING_SUCCESS";
const GET_SETTING_FAILURE = "GET_SETTING_FAILURE";

const UPDATE_SETTING = "UPDATE_SETTING";
const UPDATE_SETTING_SUCCESS = "UPDATE_SETTING_SUCCESS";
const UPDATE_SETTING_FAILURE = "UPDATE_SETTING_FAILURE";

//settingsParamsTypes
interface SettingState {
    sessionTimeout: number;
    loading: boolean;
}

//settingsActions
export const actions = {
    getSetting: () => ({
        type: GET_SETTING
    }),
    getSettingSuccess: payload => ({
        type: GET_SETTING_SUCCESS,
        payload
    }),
    getSettingFailure: (error: string) => ({
        type: GET_SETTING_FAILURE,
        error
    }),
    updateSetting: payload => ({
        type: UPDATE_SETTING,
        payload
    }),
    updateSettingSuccess: payload => ({
        type: UPDATE_SETTING_SUCCESS,
        payload
    }),
    updateSettingFailure: (error: string) => ({
        type: UPDATE_SETTING_FAILURE,
        error
    })
};

//settingsReducer
export function settingReducer(
    state: SettingState = {
        sessionTimeout: 0,
        loading: false
    },
    action
): SettingState {
    switch (action.type) {
        case GET_SETTING:
            return {
                ...state,
                loading: true
            };
        case UPDATE_SETTING_SUCCESS:
        case GET_SETTING_SUCCESS:
            const { sessionTimeout } = action.payload;
            return {
                ...state,
                sessionTimeout,
                loading: false
            };

        case LOGOUT_SUCCESS:
            return {
                ...state
            };
        default:
            return state;
    }
}

//settingsAPI
export const api = {
    getSetting: async () => {
        return await service.get(`/settings`);
    },
    updateSetting: async payload => {
        const data = {
            data: payload
        };
        return await service.put(`/settings`, data);
    }
};

//settingsSaga
function* getSettingFunc() {
    try {
        const res = yield call(api.getSetting);
        if (res) {
            yield put({ type: GET_SETTING_SUCCESS, payload: res.data });
        }
    } catch (e) {}
}

function* updateSettingFunc(action) {
    const { payload } = action;
    try {
        const res = yield call(api.updateSetting, payload);
        if (res) {
            yield put({ type: UPDATE_SETTING_SUCCESS, payload: res.data });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyUpdated",
                btnName: "ok"
            });
        }
    } catch (e) {
        Popup.alertPopup({
            className: "success",
            message: "AlertMessageFailedToUpdated",
            btnName: "ok"
        });
    }
}

export function* settingSaga() {
    yield takeLatest(GET_SETTING, getSettingFunc);
    yield takeLatest(UPDATE_SETTING, updateSettingFunc);
}
