import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { BlockInfo } from "components";
import { actions as trackerActions } from "modules/tracker";

const mapStateToProps = state => ({
    loading: state.tracker.blockLoading,
    block: state.tracker.block
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getBlock: payload => dispatch(trackerActions.getBlock(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(BlockInfo);
