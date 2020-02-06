import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { get } from '@ember/object';

export default Route.extend({
  repo: service('repository/service'),
  chainRepo: service('repository/discovery-chain'),
  settings: service('settings'),
  queryParams: {
    s: {
      as: 'filter',
      replace: true,
    },
  },
  model: function(params) {
    const dc = this.modelFor('dc').dc.Name;
    const nspace = this.modelFor('nspace').nspace.substr(1);
    return hash({
      item: this.repo.findBySlug(params.name, dc, nspace),
      urls: this.settings.findBySlug('urls'),
      dc: dc,
    }).then(model => {
      return hash({
        chain: ['connect-proxy', 'mesh-gateway'].includes(get(model, 'item.Service.Kind'))
          ? null
          : this.chainRepo.findBySlug(params.name, dc, nspace),
        ...model,
      });
    });
  },
  setupController: function(controller, model) {
    controller.setProperties(model);
  },
});
