import React, { Component } from "react"
import { Link } from "react-router-dom";

class Restrict extends Component {
    render() {
        return (
            <>
                {this.props.user.verified ? this.props.children : this.props.showWarning &&
        <div class="box">
            <div class="box-header">
                <h4 class="box-title">Verify your identity</h4>
            </div>
            <div class="box-body">
                <div class="row">
                    <div class="col-md-12">
                        <p>
                            <span class="font-weight-bold">Please verify your identity to access this feature.</span>
                            <br/>
                            <small>This is necessary to avoid abuse of the platform and ensure a safe platform for creators and the public.</small>
                        </p>
                        <p>
                            <span class="font-weight-bold">You can verify your identity by:</span>
                        </p>
                        <ul>
                            <li>
                                <span class="font-weight-bold">Instant - Making a one time payment of USD$1.00 or equiv ZWL</span>
                            </li>
                            <li>
                                <span class="font-weight-bold">Instant - Entering an invite code you received.</span>
                            </li>
                            <li>
                                <span class="font-weight-bold">72 Hours - Uploading a copy of your ID</span>
                            </li>
                        </ul>
                        <Link to="/creator/verify-identity" class="btn btn-block btn-success" ><i class="fa fa-user"></i> Verify identity</Link>
                    </div>
                </div>
            </div>
        </div>
                }
            </>
        )
    }
}

export default Restrict;