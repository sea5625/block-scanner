import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";

const GET_SYMPTOM = "GET_SYMPTOM";
const GET_SYMPTOM_SUCCESS = "GET_SYMPTOM_SUCCESS";
const GET_SYMPTOM_FAILURE = "GET_SYMPTOM_FAILURE";

//symptomParamsTypes
interface SymptomState {
    symptom: SymptomData;
    loading: boolean;
}

export interface SymptomData {
    data: {
        channel: string;
        msg: string;
        symptom: string;
        timeStamp: string;
    }[];
    total: number;
}

//symptomActions
export const actions = {
    getSymptom: payload => ({
        type: GET_SYMPTOM,
        payload
    }),
    getSymptomSuccess: (payload: SymptomState) => ({
        type: GET_SYMPTOM_SUCCESS,
        payload
    }),
    getSymptomFailure: (error: string) => ({
        type: GET_SYMPTOM_FAILURE,
        error
    })
};

//symptomReducer
export function symptomReducer(
    state: SymptomState = {
        symptom: { data: [], total: 0 },
        loading: true
    },
    action
): SymptomState {
    switch (action.type) {
        case GET_SYMPTOM:
            return {
                ...state,
                loading: true
            };
        case GET_SYMPTOM_SUCCESS:
            const { symptomData: symptom } = action;
            return {
                ...state,
                symptom,
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

//symptomAPI
export const api = {
    getSymptom: async payload => {
        const { offset, limit, from, to } = payload;
        let url = `/symptom?limit=${limit}&offset=${offset}`;
        if (from && to) {
            url = url + `&from=${from}&to=${to}`;
        }
        return await service.get(url);
    }
};

//symptomSaga
function* getSymptomFunc(action) {
    try {
        const { payload } = action;
        const res: SymptomData = yield call(api.getSymptom, payload);
        if (res) {
            yield put({ type: GET_SYMPTOM_SUCCESS, symptomData: res });
        }
    } catch (e) {
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

export function* symptomSaga() {
    yield takeLatest(GET_SYMPTOM, getSymptomFunc);
}
