import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";
import { Popup } from "components";

const GET_NODES = "GET_NODES";
const GET_NODES_SUCCESS = "GET_NODES_SUCCESS";
const GET_NODES_FAILURE = "GET_NODES_FAILURE";

const GET_ALL_NODES = "GET_ALL_NODES";
const GET_ALL_NODES_SUCCESS = "GET_ALL_NODES_SUCCESS";
const GET_ALL_NODES_FAILURE = "GET_ALL_NODES_FAILURE";

const CREATE_NODE = "CREATE_NODE";
const CREATE_NODE_SUCCESS = "CREATE_NODE_SUCCESS";
const CREATE_NODE_FAILURE = "CREATE_NODE_FAILURE";

const UPDATE_NODE = "UPDATE_NODE";
const UPDATE_NODE_SUCCESS = "UPDATE_NODE_SUCCESS";
const UPDATE_NODE_FAILURE = "UPDATE_NODE_FAILURE";

const DELETE_NODE = "DELETE_NODE";
const DELETE_NODE_SUCCESS = "DELETE_NODE_SUCCESS";
const DELETE_NODE_FAILURE = "DELETE_NODE_FAILURE";

//nodesParamsTypes
interface NodesState {
    allNodes: NodesData;
    nodes: NodesData;
    loading: boolean;
}

export interface NodesData {
    data: NodeData[];
    total: number;
}

interface NodeData {
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

//nodesActions
export const actions = {
    getAllNodes: () => ({
        type: GET_ALL_NODES
    }),
    getAllNodesSuccess: (payload: NodesData) => ({
        type: GET_ALL_NODES_SUCCESS,
        payload
    }),
    getAllNodesFailure: (error: string) => ({
        type: GET_ALL_NODES_FAILURE,
        error
    }),
    getNodes: (channelId: string) => ({
        type: GET_NODES,
        channelId
    }),
    getNodesSuccess: (payload: NodesData) => ({
        type: GET_NODES_SUCCESS,
        payload
    }),
    getNodesFailure: (error: string) => ({
        type: GET_NODES_FAILURE,
        error
    }),
    createNode: payload => ({
        type: CREATE_NODE,
        payload
    }),
    createNodeSuccess: (payload: NodesData) => ({
        type: CREATE_NODE_SUCCESS,
        payload
    }),
    createNodeFailure: (error: string) => ({
        type: CREATE_NODE_FAILURE,
        error
    }),
    updateNode: payload => ({
        type: UPDATE_NODE,
        payload
    }),
    updateNodeSuccess: (payload: NodesData) => ({
        type: UPDATE_NODE_SUCCESS,
        payload
    }),
    updateNodeFailure: (error: string) => ({
        type: UPDATE_NODE_FAILURE,
        error
    }),
    deleteNode: payload => ({
        type: DELETE_NODE,
        payload
    }),
    deleteNodeSuccess: (payload: NodesData) => ({
        type: DELETE_NODE_SUCCESS,
        payload
    }),
    deleteNodeFailure: (error: string) => ({
        type: DELETE_NODE_FAILURE,
        error
    })
};

//nodesReducer
export function nodesReducer(
    state: NodesState = {
        allNodes: { data: [], total: 0 },
        nodes: { data: [], total: 0 },
        loading: true
    },
    action
): NodesState {
    switch (action.type) {
        case GET_NODES:
        case GET_ALL_NODES:
            return {
                ...state,
                loading: true
            };
        case GET_ALL_NODES_SUCCESS:
            const { allNodesData: allNodes } = action;
            return {
                ...state,
                allNodes,
                loading: false
            };
        case GET_NODES_SUCCESS:
            const { nodesData: nodes } = action;
            return {
                ...state,
                nodes,
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

//nodesAPI
export const api = {
    getAllNodes: async () => {
        return await service.get(`/nodes`);
    },
    getNodes: async channelId => {
        return await service.get(`/nodes?channel=${channelId}`);
    },
    createNode: async payload => {
        const data = {
            data: payload
        };
        return await service.post(`/nodes`, data);
    },
    updateNode: async payload => {
        const { id, ip, name } = payload;
        const data = {
            data: {
                ip,
                name
            }
        };
        return await service.put(`/nodes/${id}`, data);
    },
    deleteNode: async payload => {
        const { id } = payload;
        return await service.delete(`/nodes/${id}`);
    }
};

//nodesSaga
function* getAllNodesFunc() {
    try {
        const res: NodesData = yield call(api.getAllNodes);
        if (res) {
            yield put({ type: GET_ALL_NODES_SUCCESS, allNodesData: res });
        }
    } catch (e) {
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}
function* getNodesFunc(action) {
    try {
        const { channelId } = action;
        const res: NodesData = yield call(api.getNodes, channelId);
        if (res) {
            yield put({ type: GET_NODES_SUCCESS, nodesData: res });
        }
    } catch (e) {
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

function* createNodeFunc(action) {
    try {
        const { payload } = action;
        const res: NodeData = yield call(api.createNode, payload);
        if (res) {
            yield put({ type: GET_ALL_NODES });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyCreate"
            });
        }
    } catch (e) {
        if (e.status === 409) {
            Popup.alertPopup({
                message: "ErrorDuplicatedNodeName",
                className: "alert",
                btnName: "ok"
            });
        }
        if (e.status === 406) {
            Popup.alertPopup({
                message: "ErrorUsedUnsupportedContentType",
                className: "alert",
                btnName: "ok"
            });
        }
        Popup.alertPopup({
            message: "ErrorFailToInsertNode",
            className: "alert",
            btnName: "ok"
        });
    }
}

function* updateNodeFunc(action) {
    try {
        const { payload } = action;
        const res: NodeData = yield call(api.updateNode, payload);
        if (res) {
            yield put({ type: GET_ALL_NODES });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyUpdated"
            });
        }
    } catch (e) {
        if (e.status === 409) {
            Popup.alertPopup({
                message: "ErrorDuplicatedNodeName",
                className: "alert",
                btnName: "ok"
            });
            return;
        }
        if (e.status === 406) {
            Popup.alertPopup({
                message: "ErrorUsedUnsupportedContentType",
                className: "alert",
                btnName: "ok"
            });
            return;
        }
        Popup.alertPopup({
            message: "ErrorFailTodUpdateNode",
            className: "alert",
            btnName: "ok"
        });
    }
}

function* deleteNodeFunc(action) {
    try {
        const { payload } = action;
        const res: NodeData = yield call(api.deleteNode, payload);
        if (res) {
            yield put({ type: GET_ALL_NODES });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyDeleted"
            });
        }
    } catch (e) {
        Popup.alertPopup({
            message: "ErrorFailToDeleteNode",
            className: "alert",
            btnName: "ok"
        });
    }
}

export function* nodesSaga() {
    yield takeLatest(GET_ALL_NODES, getAllNodesFunc);
    yield takeLatest(GET_NODES, getNodesFunc);
    yield takeLatest(CREATE_NODE, createNodeFunc);
    yield takeLatest(UPDATE_NODE, updateNodeFunc);
    yield takeLatest(DELETE_NODE, deleteNodeFunc);
}
