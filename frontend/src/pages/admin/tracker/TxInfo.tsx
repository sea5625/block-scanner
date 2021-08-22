import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { TxInfo } from "components";
import { actions as trackerActions } from "modules/tracker";

const mapStateToProps = state => ({
    loading: state.tracker.txLoading,
    tx: state.tracker.tx
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getTx: payload => dispatch(trackerActions.getTx(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(TxInfo);
