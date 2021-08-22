import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { Popup } from "../components";

const GET_USER = "GET_USER";
const GET_USER_SUCCESS = "GET_USER_SUCCESS";
const GET_USER_FAILURE = "GET_USER_FAILURE";

const GET_USER_LIST = "GET_USER_LIST";
const GET_USER_LIST_SUCCESS = "GET_USER_LIST_SUCCESS";
const GET_USER_LIST_FAILURE = "GET_USER_LIST_FAILURE";

const PUT_USER = "PUT_USER_LIST";
const PUT_USER_SUCCESS = "PUT_USER_SUCCESS";
const PUT_USER_FAILURE = "PUT_USER_FAILURE";

const CREATE_USER = "CREATE_USER";
const CREATE_USER_SUCCESS = "CREATE_USER_SUCCESS";

const DELETE_USER = "DELETE_USER";
const DELETE_USER_SUCCESS = "CREATE_USER_SUCCESS";

//usersParamsTypes
interface UsersState {
    user?: UserData;
    userList?: UserListData;
    error?: any;
    userLoading: boolean;
    userListLoading: boolean;
}

interface UserPayload {
    id: string;
}

export interface UserData {
    id: string;
    channels: Channels[];
    email: string;
    firstName: string;
    lastName: string;
    permissionToAccess: string[];
    phoneNumber: string;
    userId: string;
    userType: string;
}

interface PutUserPayload {
    id: string;
    channels: Channels[];
    email: string;
    firstName: string;
    lastName: string;
    newPassword: string;
    permissionToAccess: string[];
    phoneNumber: string;
    userId: string;
    userType: string;
}

interface CreateUserPayload {
    id: string;
    channels: Channels[];
    email: string;
    firstName: string;
    lastName: string;
    newPassword: string;
    permissionToAccess: string[];
    phoneNumber: string;
    userId: string;
    userType: string;
}

export interface UserListData {
    data: UserData[];
    total: number;
}

interface Channels {
    id: string;
    name: string;
}

//usersActions
export const actions = {
    getUser: (payload: UserPayload) => ({
        type: GET_USER,
        payload
    }),
    getUserSuccess: (payload: UserData) => ({
        type: GET_USER_SUCCESS,
        payload
    }),
    getUserFailure: (error: string) => ({
        type: GET_USER_FAILURE,
        error
    }),
    getUserList: () => ({
        type: GET_USER_LIST
    }),
    getUserListSuccess: (data: UserListData) => ({
        type: GET_USER_LIST_SUCCESS,
        data
    }),
    getUserListFailure: (error: string) => ({
        type: GET_USER_LIST_FAILURE,
        error
    }),
    putUser: (payload: PutUserPayload) => ({
        type: PUT_USER,
        payload
    }),

    createUser: (payload: CreateUserPayload) => ({
        type: CREATE_USER,
        payload
    }),

    deleteUser: (payload: UserPayload) => ({
        type: DELETE_USER,
        payload
    })
};

//usersReducer
export function usersReducer(
    state: UsersState = {
        userLoading: true,
        userListLoading: true
    },
    action
): UsersState {
    switch (action.type) {
        case GET_USER_SUCCESS:
            const { userData: user } = action;
            return {
                ...state,
                user,
                userLoading: false
            };
        case GET_USER_LIST_SUCCESS:
            const { userListData: userList } = action;
            return {
                ...state,
                userList,
                userListLoading: false
            };
        default:
            return state;
    }
}

//usersAPI
export const api = {
    getUser: async (payload: UserPayload) => {
        const { id } = payload;
        const { data } = await service.get(`/users/${id}`);
        return data;
    },
    putUser: async (payload: PutUserPayload) => {
        const _payload = {
            data: { ...payload }
        };
        delete _payload.data["id"];
        return await service.put(`users/${payload.id}`, _payload);
    },
    createUser: async (payload: CreateUserPayload) => {
        const _payload = {
            data: { ...payload }
        };
        return await service.post("/users", _payload);
    },
    getUserList: async () => {
        return await service.get("/users");
    },
    deleteUser: async (payload: UserPayload) => {
        const { id } = payload;
        return await service.delete(`/users/${id}`);
    }
};

//usersSaga
function* getUserFunc(action) {
    try {
        const payload: UserPayload = action.payload;
        const res: UserData = yield call(api.getUser, payload);
        if (res) {
            yield put({ type: GET_USER_SUCCESS, userData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

function* getUserListFunc() {
    try {
        const res: UserListData = yield call(api.getUserList);
        if (res) {
            yield put({ type: GET_USER_LIST_SUCCESS, userListData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

function* putUserFunc(action) {
    try {
        const payload: PutUserPayload = action.payload;
        const res = yield call(api.putUser, payload);
        if (res) {
            yield put({ type: PUT_USER_SUCCESS, userData: res });
            yield put({ type: GET_USER, payload });
            if (res.data.userType === "USER_ADMIN") {
                yield put({ type: GET_USER_LIST });
            }
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyUpdated"
            });
        }
    } catch (e) {
        console.log(e);
        if (e.data.errors[0].internalMessage === "Fail to update the user.") {
            Popup.alertPopup({
                className: "alert",
                message: "AlertMessageSuccessfullyUpdated"
            });
        }
    }
}

function* createUserFunc(action) {
    try {
        const payload: CreateUserPayload = action.payload;
        const res: UserData = yield call(api.createUser, payload);
        if (res) {
            yield put({ type: CREATE_USER_SUCCESS, userData: res });
            yield put({ type: GET_USER_LIST });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyCreate"
            });
        }
    } catch (e) {
        console.log(e);
        if (e.data.errors[0].internalMessage === "Fail to insert the user.") {
            Popup.alertPopup({
                className: "alert",
                message: "AlertMessageFailedToCreated"
            });
        }
        if (e.data.errors[0].internalMessage === "Duplicated user ID.") {
            Popup.alertPopup({
                className: "alert",
                message: "ErrorDuplicatedUserID"
            });
        }
    }
}

function* deleteUserFunc(action) {
    try {
        const payload: UserPayload = action.payload;
        const res: UserData = yield call(api.deleteUser, payload);
        if (res) {
            yield put({ type: DELETE_USER_SUCCESS, userData: res });
            yield put({ type: GET_USER_LIST });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyDeleted"
            });
        }
    } catch (e) {
        console.log(e);
        if (e.data.errors[0].internalMessage === "No user in DB.") {
            Popup.alertPopup({
                className: "alert",
                message: "AlertMessageFailedToDeleted"
            });
        }
    }
}

export function* usersSaga() {
    yield takeLatest(GET_USER, getUserFunc);
    yield takeLatest(GET_USER_LIST, getUserListFunc);
    yield takeLatest(CREATE_USER, createUserFunc);
    yield takeLatest(PUT_USER, putUserFunc);
    yield takeLatest(DELETE_USER, deleteUserFunc);
}
