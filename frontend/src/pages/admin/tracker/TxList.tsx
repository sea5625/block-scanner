import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { TxList } from "components";
import { actions as trackerActions } from "modules/tracker";

const mapStateToProps = state => ({
    txsLoading: state.tracker.txsLoading,
    txs: state.tracker.txs
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getTxs: payload => dispatch(trackerActions.getTxs(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(TxList);
