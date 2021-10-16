import Route from 'consul-ui/routing/route';

export default class UpstreamsRoute extends Route {
  queryParams = {
    sortBy: 'sort',
    instance: 'instance',
    searchproperty: {
      as: 'searchproperty',
      empty: [['Name', 'Tags']],
    },
    search: {
      as: 'filter',
      replace: true,
    },
  };
}
