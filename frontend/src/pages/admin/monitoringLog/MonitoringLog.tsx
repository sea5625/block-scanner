import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { MonitoringLog } from "components";
import { actions as symptomActions } from "modules/symptom";

const mapStateToProps = state => {
    return {
        symptom: state.symptom.symptom,
        loading: state.symptom.loading
    };
};

const mapDispatchToProps = (dispatch: Dispatch<Action>) => {
    return {
        getSymptom: payload => dispatch(symptomActions.getSymptom(payload))
    };
};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(MonitoringLog);
