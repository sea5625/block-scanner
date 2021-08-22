import { call, put, takeLatest } from "redux-saga/effects";
import service from "../helper/service";

//resourceActionTypes
const GET_RESOURCE = "GET_RESOURCE";
const GET_RESOURCE_SUCCESS = "GET_RESOURCE_SUCCESS";
const GET_RESOURCE_FAILURE = "GET_RESOURCE_FAILURE";

interface ResorceState {
    loading: boolean;
    resource: {
        data: string;
    };
}

interface ResorcePayload {
    id: string;
}

interface ResorceData {
    data: string;
}

//resourceActions
export const actions = {
    getResource: (payload: ResorcePayload) => ({
        type: GET_RESOURCE,
        payload
    }),
    getResourceSuccess: () => ({
        type: GET_RESOURCE_SUCCESS
    }),
    getResourceFailure: (error: string) => ({
        type: GET_RESOURCE_FAILURE,
        error
    })
};

//resourceReducer
export function resourceReducer(
    state: ResorceState = {
        loading: true,
        resource: {
            data: ""
        }
    },
    action
): ResorceState {
    switch (action.type) {
        case GET_RESOURCE_SUCCESS:
            const { resourceData: resource } = action;
            return {
                ...state,
                resource,
                loading: false
            };
        default:
            return state;
    }
}

//resourceAPI
const api = {
    getResorce: async (payload: ResorcePayload) => {
        const { id } = payload;
        return await service.get(`/resources/${id}`);
    }
};

//resourceSaga
function* getResorceFunc(action) {
    try {
        const payload: ResorcePayload = action.payload;
        const res: ResorceData = yield call(api.getResorce, payload);
        if (res) {
            yield put({ type: GET_RESOURCE_SUCCESS, resourceData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

export function* resourceSaga() {
    yield takeLatest(GET_RESOURCE, getResorceFunc);
}
