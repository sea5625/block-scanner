import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";
import { Popup } from "components";

const GET_CHANNEL = "GET_CHANNEL";
const GET_CHANNEL_SUCCESS = "GET_CHANNEL_SUCCESS";
const GET_CHANNEL_FAILURE = "GET_CHANNEL_FAILURE";

const GET_CHANNELS = "GET_CHANNELS";
const GET_CHANNELS_SUCCESS = "GET_CHANNELS_SUCCESS";
const GET_CHANNELS_FAILURE = "GET_CHANNELS_FAILURE";

const REFRESH_CHANNELS = "REFRESH_CHANNELS";
const REFRESH_CHANNELS_SUCCESS = "REFRESH_CHANNELS_SUCCESS";
const REFRESH_CHANNELS_FAILURE = "REFRESH_CHANNELS_FAILURE";

const CREATE_CHANNEL = "CREATE_CHANNEL";
const CREATE_CHANNEL_SUCCESS = "CREATE_CHANNEL_SUCCESS";
const CREATE_CHANNEL_FAILURE = "CREATE_CHANNEL_FAILURE";

const UPDATE_CHANNEL = "UPDATE_CHANNEL";
const UPDATE_CHANNEL_SUCCESS = "UPDATE_CHANNEL_SUCCESS";
const UPDATE_CHANNEL_FAILURE = "UPDATE_CHANNEL_FAILURE";

const DELETE_CHANNEL = "DELETE_CHANNEL";
const DELETE_CHANNEL_SUCCESS = "DELETE_CHANNEL_SUCCESS";
const DELETE_CHANNEL_FAILURE = "DELETE_CHANNEL_FAILURE";

//channelsParamsTypes
interface ChannelsState {
    channels: ChannelsData;
    channel: ChannelData;
    loading: boolean;
}

export interface ChannelPayload {
    id: string;
}

export interface ChannelsData {
    data: Channel[];
    total: number;
}

export interface ChannelData {
    data: Channel;
}

interface Channel {
    blockHeight: number;
    countOfTX: number;
    id: string;
    name: string;
    nodes: Nodes[];
    responseTimeInSec: number;
    status: number;
    total: number;
}

interface Nodes {
    blockHeight: number;
    countOfTX: number;
    countOfUnconfirmedTX: number;
    id: string;
    ip: string;
    isLeader: number;
    name: string;
    responseTimeInSec: number;
    status: number;
    timeStamp: string;
}

//channelsActions
export const actions = {
    getChannels: () => ({
        type: GET_CHANNELS
    }),
    getChannelsSuccess: (payload: ChannelsData) => ({
        type: GET_CHANNELS_SUCCESS,
        payload
    }),
    getChannelsFailure: (error: string) => ({
        type: GET_CHANNELS_FAILURE,
        error
    }),
    refreshChannels: () => ({
        type: REFRESH_CHANNELS
    }),
    refreshChannelsSuccess: (payload: ChannelsData) => ({
        type: REFRESH_CHANNELS_SUCCESS,
        payload
    }),
    refreshChannelsFailure: (error: string) => ({
        type: REFRESH_CHANNELS_FAILURE,
        error
    }),
    getChannel: (payload: ChannelPayload) => ({
        type: GET_CHANNEL,
        payload
    }),
    getChannelSuccess: (payload: ChannelData) => ({
        type: GET_CHANNEL_SUCCESS,
        payload
    }),
    getChannelFailure: (error: string) => ({
        type: GET_CHANNEL_FAILURE,
        error
    }),
    createChannel: payload => ({
        type: CREATE_CHANNEL,
        payload
    }),
    createChannelSuccess: (payload: ChannelData) => ({
        type: CREATE_CHANNEL_SUCCESS,
        payload
    }),
    createChannelFailure: (error: string) => ({
        type: CREATE_CHANNEL_FAILURE,
        error
    }),
    updateChannel: payload => ({
        type: UPDATE_CHANNEL,
        payload
    }),
    updateChannelSuccess: (payload: ChannelData) => ({
        type: UPDATE_CHANNEL_SUCCESS,
        payload
    }),
    updateChannelFailure: (error: string) => ({
        type: UPDATE_CHANNEL_FAILURE,
        error
    }),
    deleteChannel: payload => ({
        type: DELETE_CHANNEL,
        payload
    }),
    deleteChannelSuccess: (payload: ChannelData) => ({
        type: DELETE_CHANNEL_SUCCESS,
        payload
    }),
    deleteChannelFailure: (error: string) => ({
        type: DELETE_CHANNEL_FAILURE,
        error
    })
};

//channelsReducer
export function channelsReducer(
    state: ChannelsState = {
        channels: { data: [], total: 0 },
        channel: {
            data: {
                blockHeight: 0,
                countOfTX: 0,
                id: "",
                name: "",
                nodes: [],
                responseTimeInSec: 0,
                status: 0,
                total: 0
            }
        },
        loading: false
    },
    action
): ChannelsState {
    switch (action.type) {
        case GET_CHANNELS:
            return {
                ...state,
                loading: true
            };
        case GET_CHANNELS_SUCCESS:
        case REFRESH_CHANNELS_SUCCESS:
            const { channelsData: channels } = action;
            return {
                ...state,
                channels,
                loading: false
            };
        case GET_CHANNEL_SUCCESS:
            const { channelData: channel } = action;
            return {
                ...state,
                channel,
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

//channelsAPI
export const api = {
    getChannel: async payload => {
        const { id } = payload;
        return await service.get(`/channels/${id}`);
    },
    getChannels: async () => {
        return await service.get(`/channels`);
    },
    createChannel: async payload => {
        return await service.post(`/channels`, payload);
    },
    updateChannel: async payload => {
        const data = {
            data: payload.data
        };
        return await service.put(`/channels/${payload.id}`, data);
    },
    deleteChannel: async payload => {
        const { id } = payload;
        return await service.delete(`/channels/${id}`);
    }
};

//channelsSaga
function* getChannelFunc(action) {
    try {
        const { payload } = action;
        const res: ChannelData = yield call(api.getChannel, payload);
        if (res) {
            yield put({ type: GET_CHANNEL_SUCCESS, channelData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

function* getChannelsFunc() {
    try {
        const res: ChannelsData = yield call(api.getChannels);
        if (res) {
            yield put({ type: GET_CHANNELS_SUCCESS, channelsData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

function* refreshChannelsFunc() {
    try {
        const res: ChannelsData = yield call(api.getChannels);
        if (res) {
            yield put({ type: REFRESH_CHANNELS_SUCCESS, channelsData: res });
        }
    } catch (e) {
        console.log(e);
    }
}

function* createChannelFunc(action) {
    const { payload } = action;
    try {
        const res: ChannelData = yield call(api.createChannel, payload);
        if (res) {
            yield put({ type: GET_CHANNELS });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyCreate"
            });
        }
    } catch (e) {
        if (e.status === 409) {
            Popup.alertPopup({
                message: "ErrorDuplicatedChannelName",
                className: "alert",
                btnName: "ok"
            });
        }
    }
}

function* updateChannelFunc(action) {
    const { payload } = action;
    try {
        const res: ChannelData = yield call(api.updateChannel, payload);
        if (res) {
            yield put({ type: GET_CHANNELS });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyUpdated"
            });
        }
    } catch (e) {
        Popup.alertPopup({
            message: "ErrorFailToUpdateChannel",
            className: "alert",
            btnName: "ok"
        });
    }
}

function* deleteChannelFunc(action) {
    const { payload } = action;
    try {
        const res: ChannelData = yield call(api.deleteChannel, payload);
        if (res) {
            yield put({ type: GET_CHANNELS });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyDeleted"
            });
        }
    } catch (e) {
        Popup.alertPopup({
            message: "ErrorFailToDeleteChannel",
            className: "alert",
            btnName: "ok"
        });
    }
}

export function* channelsSaga() {
    yield takeLatest(GET_CHANNEL, getChannelFunc);
    yield takeLatest(GET_CHANNELS, getChannelsFunc);
    yield takeLatest(REFRESH_CHANNELS, refreshChannelsFunc);
    yield takeLatest(CREATE_CHANNEL, createChannelFunc);
    yield takeLatest(UPDATE_CHANNEL, updateChannelFunc);
    yield takeLatest(DELETE_CHANNEL, deleteChannelFunc);
}
