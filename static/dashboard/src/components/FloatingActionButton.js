import {Component} from "react"
import { Switch, Route, Link } from "react-router-dom";

class FloatingActionButton extends Component {

    render() {
        return (
            <Switch>
                {this.props.routes.map(route => {
                    return <Route exact path={route.route} >
                        <Link to={route.targetRoute} class="floating-action-btn"  >
                            <span class="back-top"><i class={"fa " + route.targetIcon}></i></span>
                        </Link>
                    </Route>
                })}
                {(this.props.user && this.props.user.verified) ? this.props.routes.map(route => {
                    return <Route exact path={route.targetRoute} >
                        <Link to={route.route} class="floating-action-btn"  >
                            <span class="back-top"><i class={"fa " + route.icon}></i></span>
                        </Link>
                    </Route>
                }) : <></>}
            </Switch>
            )
    }
}

export default FloatingActionButton;