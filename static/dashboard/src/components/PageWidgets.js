import React, {Component} from "react"
import { Link } from "react-router-dom"
import v1 from "../api/v1"


class PageWidgets extends Component {

    constructor(props) {
        super(props)

        this.state = {
            supporter_counter: this.props.user.page.supporter_counter,
            donate_btn: this.props.user.page.donate_btn,
            subscriber_section: this.props.user.page.subscriber_section,
            socialmedia_info: this.props.user.page.socialmedia_info,
        }

    }

    render() {
        return (
            <>
                <div className="box">
                    <div className="box-header">
                        <h4 className="box-title">Page Layout</h4>
                    </div>
                    <div className="box-body" >
                        <p>Choose which content to display on your landing page.</p>
                        <div className="form-group">
                            <div className="c-inputs-stacked">
                                <input type="checkbox" id="checkbox_347" checked={this.state.supporter_counter} />
                                <label for="checkbox_347" className="block">Show number of supporters</label>

                                <input type="checkbox" id="checkbox_123" checked={this.state.donate_btn} />
                                <label for="checkbox_123" className="block">Show buy me a {this.props.user.page.donation_item} button</label>

                                <input type="checkbox" id="checkbox_234" checked={this.state.subscriber_section} disabled={!this.props.user.memberships} />
                                <label for="checkbox_234" className="block">Show become a {this.props.user.page.donation_item} section</label>

                                <input type="checkbox" id="checkbox_346" checked={this.state.socialmedia_info} />
                                <label for="checkbox_346" className="block">Show my social media details</label>
                                {this.props.user.memberships ? <></> : <p>You must <Link to="/creator/subscriptions">Activate subscriptions</Link> first.</p>}

                            </div>

                        </div>
                    </div>
                </div>
            </>
        )
    }
}

export default PageWidgets;