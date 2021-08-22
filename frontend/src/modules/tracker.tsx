import { call, put, takeLatest } from "redux-saga/effects";
import { history } from "store";
import service from "helper/service";
import { Popup } from "components";
import { LOGOUT_SUCCESS } from "modules/auth";

const GET_BLOCKS = "GET_BLOCKS";
const GET_BLOCKS_SUCCESS = "GET_BLOCKS_SUCCESS";
const GET_BLOCKS_FAILURE = "GET_BLOCKS_FAILURE";

const GET_TXS = "GET_TXS";
const GET_TXS_SUCCESS = "GET_TXS_SUCCESS";
const GET_TXS_FAILURE = "GET_TXS_FAILURE";

const GET_BLOCK = "GET_BLOCK";
const GET_BLOCK_SUCCESS = "GET_BLOCK_SUCCESS";
const GET_BLOCK_FAILURE = "GET_BLOCK_FAILURE";

const GET_TX = "GET_TX";
const GET_TX_SUCCESS = "GET_TX_SUCCESS";
const GET_TX_FAILURE = "GET_TX_FAILURE";

const GET_BLOCK_OR_TX_BY_SEARCH = "GET_BLOCK_OR_TX_BY_SEARCH";
const GET_BLOCK_OR_TX_BY_SEARCH_SUCCESS = "GET_BLOCK_OR_TX_BY_SEARCH_SUCCESS";
const GET_BLOCK_OR_TX_BY_SEARCH_FAILURE = "GET_BLOCK_OR_TX_BY_SEARCH_FAILURE";

//trackerParamsTypes
interface TrackerState {
    blocks: BlocksData;
    txs: TxsData;
    block: BlockData;
    tx: TxData;
    blocksLoading: boolean;
    txsLoading: boolean;
    blockLoading: boolean;
    txLoading: boolean;
}

interface BlocksPayload {
    channel: string;
    limit: number;
    offset: number;
}

export interface BlocksData {
    data: BlockData[];
    total: number;
}

interface ConfirmedTx {
    blockHeight: number;
    data: { method: string };
    from: string;
    status: string;
    timeStamp: string;
    to: string;
    txHash: string;
}

interface TxsPayload {
    channel: string;
    limit: number;
    offset: number;
    status?: string;
    blockHeight?: number;
    from?: string;
    to?: string;
    fromAddress?: string;
    toAddress?: string;
    data?: string;
}

export interface TxsData {
    data: TxData[];
    total: number;
}

interface BlockPayload {
    channel: string;
    id: string;
}

export interface BlockData {
    blockHash: string;
    blockHeight: number;
    confirmedTx: ConfirmedTx[];
    peerID: string;
    signature: string;
    timeStamp: string;
}

interface TxPayload {
    channel: string;
    txHash: string;
    status?: string;
    from?: string;
    to?: string;
    blockHeight?: string;
    data?: string;
}

export interface TxData {
    txHash: string;
    blockHeight: number;
    data: {
        method: string;
    };
    status: string;
    from: string;
    to: string;
    timeStamp: string;
}

//trackerActions
export const actions = {
    getBlocks: (payload: BlocksPayload) => ({
        type: GET_BLOCKS,
        payload
    }),
    getBlocksSuccess: (payload: BlocksPayload) => ({
        type: GET_BLOCKS_SUCCESS,
        payload
    }),
    getBlocksFailure: (error: string) => ({
        type: GET_BLOCKS_FAILURE,
        error
    }),
    getTxs: (payload: TxsPayload) => ({
        type: GET_TXS,
        payload
    }),
    getTxsSuccess: (payload: TxsData) => ({
        type: GET_TXS_SUCCESS,
        payload
    }),
    getTxsFailure: (error: string) => ({
        type: GET_TXS_FAILURE,
        error
    }),
    getBlock: (payload: BlockPayload) => ({
        type: GET_BLOCK,
        payload
    }),
    getBlockSuccess: (payload: BlockPayload) => ({
        type: GET_BLOCK_SUCCESS,
        payload
    }),
    getBlockFailure: (error: string) => ({
        type: GET_BLOCK_FAILURE,
        error
    }),
    getTx: (payload: TxPayload) => ({
        type: GET_TX,
        payload
    }),
    getTxSuccess: (payload: TxData) => ({
        type: GET_TX_SUCCESS,
        payload
    }),
    getTxFailure: (error: string) => ({
        type: GET_TX_FAILURE,
        error
    }),
    getBlockOrTxBySearch: payload => ({
        type: GET_BLOCK_OR_TX_BY_SEARCH,
        payload
    }),
    getBlockOrTxBySearchSuccess: () => ({
        type: GET_BLOCK_OR_TX_BY_SEARCH_SUCCESS
    }),
    getBlockOrTxBySearchFailure: () => ({
        type: GET_BLOCK_OR_TX_BY_SEARCH_FAILURE
    })
};

//trackerReducer
export function trackerReducer(
    state: TrackerState = {
        blocks: { data: [], total: 0 },
        txs: { data: [], total: 0 },
        block: {
            blockHash: "",
            blockHeight: 0,
            confirmedTx: [],
            peerID: "",
            signature: "",
            timeStamp: ""
        },
        tx: {
            txHash: "",
            blockHeight: 0,
            data: {
                method: ""
            },
            status: "",
            from: "",
            to: "",
            timeStamp: ""
        },
        blocksLoading: false,
        txsLoading: false,
        blockLoading: false,
        txLoading: false
    },
    action
) {
    switch (action.type) {
        case GET_BLOCKS:
            return {
                ...state,
                blocksLoading: true
            };
        case GET_BLOCKS_SUCCESS:
            const { blocksData: blocks } = action;
            return {
                ...state,
                blocks,
                blocksLoading: false
            };
        case GET_TXS:
            return {
                ...state,
                txsLoading: true
            };
        case GET_TXS_SUCCESS:
            const { txsData: txs } = action;
            return {
                ...state,
                txs,
                txsLoading: false
            };
        case GET_BLOCK:
            return {
                ...state,
                blocksLoading: true
            };
        case GET_BLOCK_SUCCESS:
            const { blockData: block } = action;
            return {
                ...state,
                block,
                blockLoading: false
            };
        case GET_TX:
            return {
                ...state,
                txLoading: true
            };
        case GET_TX_SUCCESS:
            const { txData: tx } = action;
            return {
                ...state,
                tx,
                txLoading: false
            };
        case LOGOUT_SUCCESS:
            return {
                ...state
            };
        default:
            return state;
    }
}

//trackerAPI
export const api = {
    getBlocks: async payload => {
        const { channel, limit, offset } = payload;
        return await service.get(
            `/blocks?channel=${channel}&limit=${limit}&offset=${offset}`
        );
    },
    getTxs: async payload => {
        let url = "/txs?";
        Object.keys(payload).forEach(el => {
            if (payload[el] || payload[el] === 0) {
                url = url + `${el}=${payload[el]}&`;
            }
        });
        return await service.get(url);
    },
    getBlock: async payload => {
        const { id, channel } = payload;
        return await service.get(`/blocks/${id}?channel=${channel}`);
    },
    getTx: async payload => {
        const { txHash, channel } = payload;
        return await service.get(`/txs/${txHash}?channel=${channel}`);
    }
};

//trackerSaga
function* getBlocksFunc(action) {
    try {
        const { payload } = action;
        const res: BlocksData = yield call(api.getBlocks, payload);
        if (res) {
            yield put({ type: GET_BLOCKS_SUCCESS, blocksData: res });
        }
    } catch (e) {
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

function* getTxsFunc(action) {
    try {
        const { payload } = action;
        const res: TxsData = yield call(api.getTxs, payload);
        if (res) {
            yield put({ type: GET_TXS_SUCCESS, txsData: res });
        }
    } catch (e) {
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

function* getBlockFunc(action) {
    try {
        const { payload } = action;
        const res: BlockData = yield call(api.getBlock, payload);
        if (res) {
            yield put({ type: GET_BLOCK_SUCCESS, blockData: res });
        }
    } catch (e) {
        if (
            e.data.errors[0].internalMessage ===
            "Fail to query the blocks in channel from DB.  "
        ) {
            Popup.alertPopup({
                message: "noSearchResultBlock",
                className: "alert",
                btnName: "goBack",
                callbackFunc: () => history.goBack()
            });
        }
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

function* getTxFunc(action) {
    try {
        const { payload } = action;
        const res: TxData = yield call(api.getTx, payload);
        if (res) {
            yield put({ type: GET_TX_SUCCESS, txData: res });
        }
    } catch (e) {
        if (
            e.data.errors[0].internalMessage ===
            "Fail to query the Tx in channel from DB. "
        ) {
            Popup.alertPopup({
                message: "noSearchResultTxHash",
                className: "alert",
                btnName: "goBack",
                callbackFunc: () => history.goBack()
            });
        }
        if (e.status === 401) {
            yield put({ type: LOGOUT_SUCCESS });
        }
    }
}

function* getBlockOrTxBySearch(action) {
    const { channel, name, searchValue } = action.payload;
    try {
        const blockPayload = { channel, id: searchValue };
        const blockRes: BlockData = yield call(api.getBlock, blockPayload);
        if (blockRes) {
            history.push(
                `/tracker/block_info/${name}/${channel}/${searchValue}`
            );
        }
    } catch (e) {
        try {
            const txPayload = { channel, txHash: searchValue };
            const txRes: TxData = yield call(api.getTx, txPayload);
            if (txRes) {
                history.push(
                    `/tracker/tx_info/${name}/${channel}/${searchValue}`
                );
            }
        } catch (e) {
            if (
                e.data.errors[0].internalMessage ===
                "Fail to query the Tx in channel from DB. "
            ) {
                Popup.alertPopup({
                    message: "noSearchResult",
                    className: "alert",
                    btnName: "ok"
                });
            }
        }
    }
}

export function* trackerSaga() {
    yield takeLatest(GET_BLOCKS, getBlocksFunc);
    yield takeLatest(GET_TXS, getTxsFunc);
    yield takeLatest(GET_BLOCK, getBlockFunc);
    yield takeLatest(GET_TX, getTxFunc);
    yield takeLatest(GET_BLOCK_OR_TX_BY_SEARCH, getBlockOrTxBySearch);
}
