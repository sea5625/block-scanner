import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { Dashboard } from "components";
import { actions as channelsActions } from "modules/channels";
import { actions as prometheusActions } from "modules/prometheus";
import { actions as alertingActions } from "modules/alerting";

const mapStateToProps = state => ({
    channels: state.channels.channels,
    storage: state.storage,
    prometheus: state.prometheus.prometheus,
    alerting: state.alerting.alerting
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    refreshChannels: () => dispatch(channelsActions.refreshChannels()),
    getPrometheus: () => dispatch(prometheusActions.getPrometheus()),
    getAlerting: () => dispatch(alertingActions.getAlerting())
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Dashboard);
