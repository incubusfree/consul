import Component from '@glimmer/component';
import { action } from '@ember/object';
import { trackedInLocalStorage } from 'ember-tracked-local-storage';

export default class AgentlessNotice extends Component {
  @trackedInLocalStorage({ defaulValue: 'false' }) consulNodesAgentlessNoticeDismissed;

  get isVisible() {
    const { items, filteredItems } = this.args;

    console.log(this.consulNodesAgentlessNoticeDismissed !== 'true' && items.length > filteredItems.length);
    console.log('tracked prop: ', this.consulNodesAgentlessNoticeDismissed);
    return this.consulNodesAgentlessNoticeDismissed !== 'true' && items.length > filteredItems.length;
  }

  @action
  dismissAgentlessNotice() {
    this.consulNodesAgentlessNoticeDismissed = 'true';
  }
}
