#!/usr/bin/env python3
"""生成带有详细参数的特性配置文件"""

import csv
import json
import urllib.request
from io import StringIO

BASE_URL = "https://raw.githubusercontent.com/PokeAPI/pokeapi/master/data/v2/csv"

def to_camel_case(snake_str):
    """将 snake_case 转换为 CamelCase"""
    components = snake_str.split('-')
    return ''.join(x.title() for x in components)

def fetch_csv(filename):
    """获取 CSV 数据"""
    url = f"{BASE_URL}/{filename}"
    with urllib.request.urlopen(url, timeout=60) as response:
        content = response.read().decode('utf-8')
    reader = csv.reader(StringIO(content))
    next(reader)
    return list(reader)

# 特性效果配置模板
ABILITY_CONFIGS = {
    # 伤害修正类
    1: {"trigger": "on_attack", "flinch_chance": 0.1},  # 恶臭
    3: {"trigger": "turn_end", "stat": "speed", "stages": 1},  # 加速
    22: {"trigger": "on_switch_in", "stat": "attack", "stages": -1, "target": "opponent"},  # 威吓
    26: {"trigger": "before_damage", "immune_types": ["ground"]},  # 飘浮
    37: {"trigger": "stat_calc", "stat": "attack", "multiplier": 2.0},  # 大力士
    47: {"trigger": "type_effectiveness", "types": ["fire", "ice"], "multiplier": 0.5},  # 厚脂肪
    65: {"trigger": "damage_calc", "condition": "hp_below_third", "type_boost": "grass", "multiplier": 1.5},  # 茂盛
    66: {"trigger": "damage_calc", "condition": "hp_below_third", "type_boost": "fire", "multiplier": 1.5},  # 猛火
    67: {"trigger": "damage_calc", "condition": "hp_below_third", "type_boost": "water", "multiplier": 1.5},  # 激流
    74: {"trigger": "stat_calc", "stat": "attack", "multiplier": 2.0},  # 瑜伽之力
    91: {"trigger": "damage_calc", "stab_multiplier": 2.0},  # 适应力
    97: {"trigger": "damage_calc", "critical_multiplier": 2.25},  # 狙击手
    101: {"trigger": "damage_calc", "condition": "power_le_60", "multiplier": 1.5},  # 技术高手
    136: {"trigger": "before_damage", "condition": "hp_full", "multiplier": 0.5},  # 多重鳞片
    153: {"trigger": "on_ko", "stat": "attack", "stages": 1},  # 自信过度
    173: {"trigger": "damage_calc", "move_flags": ["bite"], "multiplier": 1.5},  # 强壮之颚
    181: {"trigger": "damage_calc", "move_flags": ["contact"], "multiplier": 1.3},  # 硬爪
    
    # 状态异常类
    9: {"trigger": "on_contact", "status": "paralysis", "chance": 0.3},  # 静电
    17: {"trigger": "immunity", "status": "poison"},  # 免疫
    38: {"trigger": "on_contact", "status": "poison", "chance": 0.3},  # 毒刺
    49: {"trigger": "on_contact", "status": "burn", "chance": 0.3},  # 火焰之躯
    62: {"trigger": "status_active", "status": "any", "stat": "attack", "multiplier": 1.5},  # 毅力
    
    # 天气相关
    2: {"trigger": "on_switch_in", "weather": "rain"},  # 降雨
    33: {"trigger": "weather_active", "weather": "rain", "stat": "speed", "multiplier": 2.0},  # 悠游自如
    34: {"trigger": "weather_active", "weather": "sun", "stat": "speed", "multiplier": 2.0},  # 叶绿素
    45: {"trigger": "on_switch_in", "weather": "sandstorm"},  # 扬沙
    70: {"trigger": "on_switch_in", "weather": "sun"},  # 日照
    117: {"trigger": "on_switch_in", "weather": "hail"},  # 降雪
    146: {"trigger": "weather_active", "weather": "sandstorm", "stat": "speed", "multiplier": 2.0},  # 拨沙
    202: {"trigger": "weather_active", "weather": "hail", "stat": "speed", "multiplier": 2.0},  # 拨雪
    
    # 吸��类
    10: {"trigger": "hit_by_type", "type": "electric", "heal_percent": 0.25},  # 蓄电
    11: {"trigger": "hit_by_type", "type": "water", "heal_percent": 0.25},  # 储水
    18: {"trigger": "hit_by_type", "type": "fire", "special_attack_boost": 1},  # 引火
    31: {"trigger": "hit_by_type", "type": "electric", "stat": "special_attack", "stages": 1},  # 避雷针
    114: {"trigger": "hit_by_type", "type": "water", "stat": "special_attack", "stages": 1},  # 引水
    
    # 变化招式相关
    156: {"trigger": "hit_by_move", "category": "status", "reflect": True},  # 魔法镜
    158: {"trigger": "priority_boost", "category": "status", "boost": 1},  # 恶作剧之心
    
    # 特殊机制
    25: {"trigger": "before_damage", "condition": "not_super_effective", "immune": True},  # 神奇守护
    98: {"trigger": "immunity", "damage_types": ["indirect"]},  # 魔法防守
    168: {"trigger": "before_move", "type_change": "move_type"},  # 变幻自如
    177: {"trigger": "priority_boost", "condition": "hp_full", "type": "flying", "boost": 1},  # 疾风之翼
    236: {"trigger": "before_move", "type_change": "move_type"},  # 自由者
}

def main():
    print("正在获取特性基础数据...")
    abilities_csv = fetch_csv("abilities.csv")
    
    print("正在获取特性中文名称...")
    names_csv = fetch_csv("ability_names.csv")
    
    # 构建映射
    ability_identifiers = {}
    for row in abilities_csv:
        if len(row) >= 2:
            ability_id = int(row[0])
            identifier = row[1]
            ability_identifiers[ability_id] = identifier
    
    ability_names = {}
    for row in names_csv:
        if len(row) >= 3:
            ability_id = int(row[0])
            lang_id = int(row[1])
            name = row[2]
            if lang_id == 12:  # 简体中文
                ability_names[ability_id] = name
    
    # 加载已有的描述
    print("正在加载已有特性描述...")
    with open("assets/pokemon/abilities.json", "r", encoding="utf-8") as f:
        existing_data = json.load(f)
    
    effect_map = {a["id"]: a["effect"] for a in existing_data["abilities"]}
    
    # 构建增强配置
    abilities = []
    for ability_id in sorted(ability_names.keys()):
        identifier = ability_identifiers.get(ability_id, f"unknown_{ability_id}")
        name = ability_names[ability_id]
        effect = effect_map.get(ability_id, "")
        
        # 生成变量名
        var_name = f"Ability{to_camel_case(identifier)}"
        
        # 基础配置
        config = {
            "id": ability_id,
            "identifier": identifier,
            "var_name": var_name,
            "name": name,
            "effect": effect,
        }
        
        # 添加详细参数（如果已定义）
        if ability_id in ABILITY_CONFIGS:
            config["params"] = ABILITY_CONFIGS[ability_id]
        else:
            config["params"] = {"trigger": "not_implemented"}
        
        abilities.append(config)
    
    # 输出 JSON
    output = {
        "version": "1.0.0",
        "total": len(abilities),
        "abilities": abilities
    }
    
    with open("assets/pokemon/abilities.json", "w", encoding="utf-8") as f:
        json.dump(output, f, ensure_ascii=False, indent=2)
    
    print(f"完成！共 {len(abilities)} 个特性，已保存到 assets/pokemon/abilities.json")
    print(f"已配置参数的特性：{len(ABILITY_CONFIGS)} 个")

if __name__ == "__main__":
    main()
