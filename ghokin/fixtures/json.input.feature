@unit
Feature: Unit: ListUnits

  Scenario: I list units
    When I call the "Unit"."ListUnits" handler with:
      """
      {
        "filters": {
          "sort": {
            "keys": [
              "_key"
            ],
            "mode": 1
          }
        }
      }
      """
    Then the response should be successful
    And the response result JSON should be a superset of:
        """
        {
            "units": [
            {
                "id": "unit_key",
         "name": "unit_name",
                "symbol": "unit_symbol",
                "quantity": "unit_quantity",
                "plural": "unit_plural",
                "aliases": [
                  "unit_alias1",
               "unit_alias2"
                                         ],
                "system": "unit_system",
                "user_defined": true
            },
            {
      "id": "unit_key2",
      "name": "unit_name2",
      "symbol": "unit_symbol2",
      "quantity": "unit_quantity2",
      "plural": "unit_plural2",
                "aliases": [
                        "unit_alias3",
                        "unit_alias4"
                        ],
                        "system": "unit_system2",
                "user_defined": true
            }
          ]
        }
        """
