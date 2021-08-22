import React, { Component } from "react";
import { IconButton } from "@material-ui/core";
import { MdClose } from "react-icons/md";
import { LayerPopup } from "lib/popup";

class PopupContainer extends Component {
    componentWillUnmount() {
        document.body.style.overflow = "visible";
    }

    componentDidUpdate(prevProps, prevState) {
        document.body.style.overflow = "hidden";
    }

    componentDidMount() {
        document.body.style.overflow = "hidden";
    }

    onClickClose = () => {
        LayerPopup.hide(this.props.layerKey);
    };

    render() {
        const newProps = {
            location: this.props.location,
            history: this.props.history,
            layerKey: this.props.layerKey,
            layerCount: this.props.layerCount,
            LayerPopup: LayerPopup
        };
        const { className } = this.props.children.props;
        return (
            <React.Fragment>
                <div
                    layerkey={this.props.layerKey}
                    className={`popup-container ${className ? className : ""}`}
                >
                    <div className="dimde">
                        <div className="popup">
                            <IconButton
                                className="close-btn"
                                onClick={this.onClickClose}
                            >
                                <MdClose />
                            </IconButton>
                            {React.cloneElement(
                                this.props.children,
                                newProps,
                                this.props
                            )}
                        </div>
                    </div>
                </div>
            </React.Fragment>
        );
    }
}

export default PopupContainer;
