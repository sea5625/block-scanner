import React from "react";
import { Link } from "react-router-dom";
import { MdLaptopChromebook, MdList } from "react-icons/md";
import { DiNetbeans } from "react-icons/di";

const Total = props => {
    return (
        <div className="total">
            <div className="total-item node">
                <Link
                    to={{
                        pathname: `/node_list/${props.name}`,
                        state: {
                            channelId: props.channelId
                        }
                    }}
                >
                    <MdLaptopChromebook />
                    {props.channelTotal >= 6 ? (
                        <p className="total-title">
                            Number of <br />
                            Nodes
                        </p>
                    ) : (
                        <p className="total-title">Number of Nodes</p>
                    )}
                    <p className="total-value">{props.totalNodes}</p>
                </Link>
            </div>
            <div className="total-item block">
                <Link
                    to={{
                        pathname: `/tracker/block_list/${props.name}/${props.channelId}`
                    }}
                >
                    <DiNetbeans />
                    {props.channelTotal >= 6 ? (
                        <p className="total-title">
                            Total <br />
                            Block Height
                        </p>
                    ) : (
                        <p className="total-title">Total Block Height</p>
                    )}
                    <p className="total-value">{props.blockHeight}</p>
                </Link>
            </div>
            <div className="total-item tx">
                <Link
                    to={{
                        pathname: `/tracker/tx_list/${props.name}/${props.channelId}`
                    }}
                >
                    <MdList />
                    {props.channelTotal >= 6 ? (
                        <p className="total-title">
                            Total <br />
                            Transactions
                        </p>
                    ) : (
                        <p className="total-title">Total Transactions</p>
                    )}
                    <p className="total-value">{props.countOfTX}</p>
                </Link>
            </div>
        </div>
    );
};

export default Total;
