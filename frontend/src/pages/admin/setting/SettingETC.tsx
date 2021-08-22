import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { SettingETC } from "components";
import { actions as alertingActions } from "modules/alerting";

const mapStateToProps = state => ({
    alerting: state.alerting.alerting,
    alertingLoading: state.alerting.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getAlerting: () => dispatch(alertingActions.getAlerting()),
    updateAlerting: payload => dispatch(alertingActions.updateAlerting(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SettingETC);
